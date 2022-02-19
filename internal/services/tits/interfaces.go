package tits

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type Database interface {
	GetTits(ctx context.Context) ([]domain.Tits, error)
	CreateTits(ctx context.Context, tits domain.Tits) error
	IncreaseRating(ctx context.Context, titsID string) (int64, error)
}

type Storage interface {
	CreateImageFromFile(ctx context.Context, imageName string, filePath string) error
	CreateImageFromBytes(ctx context.Context, imageName string, imageData []byte) error
	GetImageUrl(imageID string) string
}
