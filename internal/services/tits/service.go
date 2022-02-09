package tits

import (
	"context"
	"strings"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const defaultTitsCreateTimeout = time.Second * 10

type Service struct {
	db      Database
	storage Storage

	wsChannel chan domain.WSMessage

	log *otelzap.Logger
}

func NewService(db Database, storage Storage, log *zap.Logger, wsChannel chan domain.WSMessage) *Service {
	return &Service{
		db:        db,
		storage:   storage,
		wsChannel: wsChannel,
		log:       otelzap.New(log.Named("tits_service")),
	}
}

func (s *Service) CreateTitsFromFile(ctx context.Context, filename, filePath string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTitsCreateTimeout)
	defer cancel()

	err := s.storage.CreateImageFromFile(ctx, filename, filePath)
	if err != nil {
		s.log.Ctx(ctx).Error("create tits from file:", zap.Error(err))
		return err
	}
	err = s.db.CreateTits(ctx, domain.Tits{
		ID:        strings.ReplaceAll(filename, ".jpg", ""),
		CreatedAt: time.Now().UTC(),
		Rating:    0,
	})
	if err != nil {
		s.log.Ctx(ctx).Error("create tits in db: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) GetTits(ctx context.Context) ([]domain.Tits, error) {
	tits, err := s.db.GetTits(ctx)
	if err != nil {
		s.log.Ctx(ctx).Error("get tits from db", zap.Error(err))
		return nil, err
	}

	for idx := range tits {
		s.storage.AssembleFileURL(&tits[idx])
	}

	return tits, nil
}

func (s *Service) IncreaseRating(ctx context.Context, titsID string) error {
	newRating, err := s.db.IncreaseRating(ctx, titsID)
	if err != nil {
		s.log.Ctx(ctx).Error("increase rating in db", zap.Error(err))
		return err
	}

	go s.sendNewRatingMessage(titsID, newRating)

	return nil
}

func (s *Service) sendNewRatingMessage(titsID string, newRating int64) {
	s.wsChannel <- domain.WSMessage{
		Type: domain.WSMessageTypeNewRating,
		Message: domain.WSNewRatingMessage{
			TitsID:    titsID,
			NewRating: newRating,
		},
	}
}
