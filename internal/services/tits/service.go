package tits

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const defaultTitsCreateTimeout = time.Second * 10

type Service struct {
	db           Database
	storage      Storage
	optimizerURL string

	wsChannel chan domain.WSMessage

	log *otelzap.Logger
}

func NewService(db Database, storage Storage, log *zap.Logger, wsChannel chan domain.WSMessage, optimizerURL string) *Service {
	return &Service{
		db:           db,
		storage:      storage,
		wsChannel:    wsChannel,
		optimizerURL: optimizerURL,
		log:          otelzap.New(log.Named("tits_service")),
	}
}

func (s *Service) getWebpImage(ctx context.Context, filename string) ([]byte, error) {
	httpClient := http.Client{}
	filenameSplitted := strings.Split(filename, ".")
	fileUrl := s.storage.GetImageUrl(filenameSplitted[0])
	requestURL := fmt.Sprintf("%s/optimize?size=350&format=webp&src=http://minio.images:9000%s.jpg", s.optimizerURL, fileUrl)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *Service) CreateTitsFromFile(ctx context.Context, filename, filePath string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTitsCreateTimeout)
	defer cancel()

	err := s.storage.CreateImageFromFile(ctx, filename, filePath)
	if err != nil {
		s.log.Ctx(ctx).Error("create tits from file:", zap.Error(err))
		return err
	}

	webpImage, err := s.getWebpImage(ctx, filename)
	if err != nil {
		s.log.Ctx(ctx).Error("get webp image:", zap.Error(err))
		return err
	}

	webpFilename := strings.Replace(filename, ".jpg", ".webp", 1)
	err = s.storage.CreateImageFromBytes(ctx, webpFilename, webpImage)
	if err != nil {
		s.log.Ctx(ctx).Error("create webp image:", zap.Error(err))
		return err
	}

	//err = s.db.CreateTits(ctx, domain.Tits{
	//	ID:        strings.ReplaceAll(filename, ".jpg", ""),
	//	CreatedAt: time.Now().UTC(),
	//	Rating:    0,
	//})
	//if err != nil {
	//	s.log.Ctx(ctx).Error("create tits in db: ", zap.Error(err))
	//	return err
	//}
	return nil
}

func (s *Service) GetTits(ctx context.Context) ([]domain.Tits, error) {
	tits, err := s.db.GetTits(ctx)
	if err != nil {
		s.log.Ctx(ctx).Error("get tits from db", zap.Error(err))
		return nil, err
	}

	for idx := range tits {
		imgPrefix := s.storage.GetImageUrl(tits[idx].ID)
		tits[idx].URL = fmt.Sprintf("%s.webp", imgPrefix)
		tits[idx].FullURL = fmt.Sprintf("%s.jpg", imgPrefix)
	}

	return tits, nil
}

func (s *Service) GetTop(ctx context.Context, limit int) ([]domain.Tits, error) {
	tits, err := s.db.GetTop(ctx, limit)
	if err != nil {
		s.log.Ctx(ctx).Error("get tits from db", zap.Error(err))
		return nil, err
	}

	for idx := range tits {
		imgPrefix := s.storage.GetImageUrl(tits[idx].ID)
		tits[idx].URL = fmt.Sprintf("%s.webp", imgPrefix)
		tits[idx].FullURL = fmt.Sprintf("%s.jpg", imgPrefix)
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
