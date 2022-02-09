package initiator

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const titsPath = "assets/images"

type Service struct {
	log         *otelzap.Logger
	titsService TitsService
}

func NewService(log *zap.Logger, titsService TitsService) *Service {
	return &Service{
		log:         otelzap.New(log.Named("initiator")),
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

	for idx, f := range files {
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
	}
}
