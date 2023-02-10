package initiator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.Logger
	titsService TitsService
	taskService TaskRepo
	httpClient  *http.Client
}

func NewService(log *zap.Logger, titsService TitsService, taskService TaskRepo) *Service {
	return &Service{
		log:         log.Named("initiator"),
		titsService: titsService,
		taskService: taskService,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *Service) Run() {
	s.log.Info("Tits downloader started")
	defer s.log.Info("Tits downloader stopped")
	ctx := context.Background()

	totalTasks, err := s.taskService.GetCountUnprocessedTasks(ctx)
	if err != nil {
		s.log.Error("get total tasks", zap.Error(err))
		return
	}

	guard := make(chan struct{}, 1000)
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

func (s *Service) RunFill() {
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
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			allUrls = append(allUrls, url)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory assets/urls/: %v\n", err)
		os.Exit(1)
	}

	totalFiles := len(allUrls)

	s.log.Info("Total urls", zap.Int("count", totalFiles))

	for idx := range allUrls {
		err := s.taskService.CreateTask(ctx, domain.Task{
			ID:        domain.NewID(),
			CreatedAt: time.Now(),
			Processed: false,
			Url:       allUrls[idx],
			Status:    "",
		})
		if err != nil {
			s.log.Error("create task: ", zap.Error(err))
		}
		s.log.Info("Task created", zap.Int("index", idx), zap.Int("total", totalFiles))
		continue
		//guard <- struct{}{}
		//wg.Add(1)
		//go s.work(ctx, wg, guard, idx, allUrls[idx], totalFiles)
	}
	//wg.Wait()
}

func (s *Service) work(ctx context.Context, wg *sync.WaitGroup, guard chan struct{}, idx int, url string, totalFiles int, task domain.Task) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*90)
	s.log.Info(
		"Creating new tits",
		zap.String("url", url),
		zap.Int("index", idx),
		zap.Int("total", totalFiles),
	)

	res, err := s.httpClient.Get(url)
	if err != nil {
		fmt.Printf("Error downloading image from URL %s: %v\n", url, err)
		task.Processed = true
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	if res.StatusCode != 200 {
		fmt.Printf("Error downloading image from URL %s: Status code: %d\n", url, res.StatusCode)
		task.Processed = true
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		task.Processed = true
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
	}
	_ = res.Body.Close()

	// If 'b' size less than 700kb, return
	if len(b) < 600*1024 {
		task.Processed = true
		_ = s.taskService.UpdateTask(context.Background(), task)
		cancel()
		wg.Done()
		<-guard
		return
	}

	err = s.titsService.CreateTitsFromBytes(ctx, fmt.Sprintf("%s.jpg", task.ID), b)
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
