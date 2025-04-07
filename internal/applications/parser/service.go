package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	log              *zap.Logger
	titsService      TitsService
	taskService      TaskRepo
	detectionService DetectorService
	httpClient       *http.Client
	httpProxyClient  *http.Client
}

func NewService(
	log *zap.Logger,
	titsService TitsService,
	taskService TaskRepo,
	detectionService DetectorService,
	proxyUrl string,
) *Service {
	proxyTransport := &http.Transport{
		TLSHandshakeTimeout: 30 * time.Second,
	}

	parsedProxyUrl, err := url.Parse(proxyUrl)
	if err == nil {
		proxyTransport.Proxy = http.ProxyURL(parsedProxyUrl)
	}

	return &Service{
		log:              log.Named("initiator"),
		titsService:      titsService,
		taskService:      taskService,
		detectionService: detectionService,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		httpProxyClient: &http.Client{
			Timeout:   time.Second * 60,
			Transport: proxyTransport,
		},
	}
}

func (s *Service) Run(withFill bool) {
	s.log.Info("Tits downloader started")
	defer s.log.Info("Tits downloader stopped")

	s.run(withFill)
}

func (s *Service) run(withFill bool) {
	ctx := context.Background()

	if withFill {
		s.runFillFormFiles()
	}

	totalTasks, err := s.taskService.GetCountUnprocessedTasks(ctx)
	if err != nil {
		s.log.Error("get total tasks", zap.Error(err))
		return
	}

	guard := make(chan struct{}, 50)
	wg := &sync.WaitGroup{}

	for i := 0; i <= totalTasks; i++ {
		task, err := s.taskService.GetTask(ctx)
		if err != nil {
			s.log.Error("get task", zap.Error(err))
			continue
		}
		s.log.Info("process idx", zap.Int("idx", i))
		guard <- struct{}{}
		wg.Add(1)
		go s.work(ctx, wg, guard, i, task.Url, totalTasks, task)
	}

	wg.Wait()
}

func (s *Service) runFillFormFiles() {
	s.log.Info("Tits uploader started")
	defer s.log.Info("Tits uploader stopped")
	ctx := context.Background()

	//guard := make(chan struct{}, 500)
	//wg := &sync.WaitGroup{}
	var allUrls []string

	err := filepath.Walk("assets/urls/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			s.log.Error("open file", zap.Error(err))
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			uri := scanner.Text()
			allUrls = append(allUrls, uri)
		}

		if err := scanner.Err(); err != nil {
			s.log.Error("scan file", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Error("walking directory assets/urls/", zap.Error(err))
		os.Exit(1)
	}

	totalFiles := len(allUrls)

	s.log.Info("Total urls", zap.Int("count", totalFiles))

	var tasks []domain.Task

	for idx := range allUrls {
		// create tasks in batches of 2000
		if len(tasks) >= 2000 {
			err := s.taskService.CreateTask(ctx, tasks)
			if err != nil {
				s.log.Error("create task: ", zap.Error(err))
			}
			s.log.Info("Task created", zap.Int("index", idx), zap.Int("total", totalFiles))
			tasks = []domain.Task{}
		} else {
			tasks = append(tasks, domain.Task{
				ID:        domain.NewID(),
				CreatedAt: time.Now(),
				Processed: false,
				Url:       allUrls[idx],
				Status:    "",
			})
		}
	}

	err = s.taskService.CreateTask(ctx, tasks)
	if err != nil {
		s.log.Error("create task: ", zap.Error(err))
	}
}

func (s *Service) work(ctx context.Context, wg *sync.WaitGroup, guard chan struct{}, idx int, url string, totalFiles int, task domain.Task) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*90)
	s.log.Info(
		"Creating new tits",
		zap.String("url", url),
		zap.Int("index", idx),
		zap.Int("total", totalFiles),
	)

	var res *http.Response
	var err error

	detectionResult, err := s.detectionService.Detect(ctx, url)
	if err != nil {
		s.log.Error("Failed to detect content",
			zap.String("url", url),
			zap.Int("index", idx),
			zap.Int("total", totalFiles),
			zap.Error(err),
		)
		if task.NeedRetry {
			s.log.Info("Detection failed, retrying",
				zap.String("url", url),
				zap.Int("index", idx),
				zap.Int("total", totalFiles),
			)
		} else {
			task.NeedRetry = true
			s.log.Info("Detection failed, not retrying",
				zap.String("url", url),
				zap.Int("index", idx),
				zap.Int("total", totalFiles),
			)
		}
		task.Processed = false
		task.Error = err.Error()
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	var detectionThresholds = map[domain.DetectionClass]float64{
		domain.DetectionClassFemaleGenitaliaCovered: 0.5,
		domain.DetectionClassFemaleGenitaliaExposed: 0.3,
		domain.DetectionClassMaleBreastExposed:      0.3,
		domain.DetectionClassAnusExposed:            0.3,
		domain.DetectionClassFaceMale:               0.3,
		domain.DetectionClassMaleGenitaliaExposed:   0.3,
		domain.DetectionClassAnusCovered:            0.5,
	}

	task.DetectionResult = detectionResult

	if len(detectionResult.Detections) > 0 {
		for _, detection := range detectionResult.Detections {
			if threshold, ok := detectionThresholds[detection.Class]; ok {
				if detection.Score > threshold {
					s.log.Info("Detection score is above threshold",
						zap.String("url", url),
						zap.Int("index", idx),
						zap.Int("total", totalFiles),
						zap.String("detection_class", string(detection.Class)),
						zap.Float64("score", detection.Score),
					)
					task.NeedRetry = false
					task.DetectionResult = detectionResult
					task.Processed = true
					task.Error = fmt.Sprintf("Detection score is above threshold: %s > %f", detection.Class, threshold)
					_ = s.taskService.UpdateTask(context.Background(), task)
					cancel()
					wg.Done()
					<-guard
					return
				}
			} else {
				s.log.Info("Detection class not in thresholds",
					zap.String("url", url),
					zap.Int("index", idx),
					zap.Int("total", totalFiles),
					zap.String("detection_class", string(detection.Class)),
					zap.Float64("score", detection.Score),
				)
			}
		}
	}

	if task.NeedRetry {
		res, err = s.httpProxyClient.Get(url)
	} else {
		res, err = s.httpClient.Get(url)
	}

	if err != nil {
		s.log.Error("downloading image from URL", zap.Error(err), zap.String("url", url))
		if !task.NeedRetry {
			task.NeedRetry = true
		} else {
			task.Processed = true
			task.Error = err.Error()
		}
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	if res.StatusCode != 200 {
		s.log.Error("downloading image from URL", zap.Error(err), zap.String("url", url), zap.Int("status", res.StatusCode))
		if !task.NeedRetry {
			task.NeedRetry = true
		} else {
			task.Processed = true
			task.Error = fmt.Sprintf("Status code = %d", res.StatusCode)
		}
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if !task.NeedRetry {
			task.NeedRetry = true
		} else {
			task.Processed = true
			task.Error = err.Error()
		}
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
	}
	_ = res.Body.Close()

	// If 'b' size less than 400kb, return
	if len(b) < 200*1024 {
		task.Processed = true
		task.Error = "image size less than 400kb"
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	err = s.titsService.CreateTitsFromBytes(ctx, fmt.Sprintf("%s.jpg", task.ID), b, url)
	if err != nil {
		s.log.Error("Failed to create tits",
			zap.String("url", url),
			zap.Int("index", idx),
			zap.Int("total", totalFiles),
			zap.Error(err),
		)
	}

	task.Processed = true
	_ = s.taskService.UpdateTask(context.Background(), task)
	cancel()
	wg.Done()
	<-guard
}
