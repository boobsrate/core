package detector

import (
	"context"

	"github.com/boobsrate/core/internal/clients/detection"
	"github.com/boobsrate/core/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.Logger
	client *detection.Client
}

func NewService(log *zap.Logger, client *detection.Client) *Service {
	return &Service{
		log:    log.Named("detector"),
		client: client,
	}
}

func (s *Service) Detect(ctx context.Context, url string) (domain.DetectionResult, error) {
	s.log.Info("detecting content", zap.String("url", url))

	result, err := s.client.Detect(ctx, url)
	if err != nil {
		s.log.Error("failed to detect content", zap.Error(err), zap.String("url", url))
		return domain.DetectionResult{}, err
	}

	return result, nil
}
