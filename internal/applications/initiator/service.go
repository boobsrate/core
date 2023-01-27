package initiator

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"sync"

	"go.uber.org/zap"
)

const titsPath = "assets/images"

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

	files, err := ioutil.ReadDir(titsPath)
	if err != nil {
		s.log.Error("Failed to read directory", zap.Error(err))
	}

	ctx := context.Background()
	totalFiles := len(files)

	guard := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for idx, f := range files {
		guard <- struct{}{}
		wg.Add(1)
		go s.work(ctx, wg, guard, idx, totalFiles, f)
	}
	wg.Wait()
}

func (s *Service) work(ctx context.Context, wg *sync.WaitGroup, guard chan struct{}, idx, totalFiles int, f fs.FileInfo) {
	s.log.Info(
		"Creating new tits",
		zap.String("name", f.Name()),
		zap.Int("index", idx),
		zap.Int("total", totalFiles),
	)
	err := s.titsService.CreateTitsFromFile(ctx, f.Name(), fmt.Sprintf("%s/%s", titsPath, f.Name()))
	if err != nil {
		s.log.Error("Failed to create tits",
			zap.String("name", f.Name()),
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
