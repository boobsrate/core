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

	"github.com/boobsrate/core/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.Logger
	titsService TitsService
}

func NewService(log *zap.Logger, titsService TitsService) *Service {
	return &Service{
		log:         log.Named("initiator"),
		titsService: titsService,
	}
}

func (s *Service) Run() {
	s.log.Info("Tits uploader started")
	defer s.log.Info("Tits uploader stopped")
	ctx := context.Background()

	guard := make(chan struct{}, 300)
	wg := &sync.WaitGroup{}
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
		if idx <= 6081 {
			continue
		}
		guard <- struct{}{}
		wg.Add(1)
		go s.work(ctx, wg, guard, idx, allUrls[idx], totalFiles)
	}
	wg.Wait()
}

func (s *Service) work(ctx context.Context, wg *sync.WaitGroup, guard chan struct{}, idx int, url string, totalFiles int) {
	s.log.Info(
		"Creating new tits",
		zap.String("url", url),
		zap.Int("index", idx),
		zap.Int("total", totalFiles),
	)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading image from URL %s: %v\n", url, err)
		return
	}

	if res.StatusCode != 200 {
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	err = s.titsService.CreateTitsFromBytes(ctx, fmt.Sprintf("%s.jpg", domain.NewID()), b)
	if err != nil {
		s.log.Error("Failed to create tits",
			zap.String("url", url),
			zap.Int("index", idx),
			zap.Int("total", totalFiles),
			zap.Error(err),
		)
	}
	defer func() {
		wg.Done()
		<-guard
	}()
}
