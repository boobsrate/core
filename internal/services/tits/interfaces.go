package tits

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type Database interface {
	GetTits(ctx context.Context) ([]domain.Tits, error)
	GetTop(ctx context.Context, limit int, abyss bool) ([]domain.Tits, error)
	CreateTits(ctx context.Context, tits domain.Tits) error
	IncreaseRating(ctx context.Context, titsID string) (int64, error)
	Report(ctx context.Context, titsID string) error
	GetReportsCount(ctx context.Context, titsID string) (int, error)
	MoveToAbyss(ctx context.Context, titsID string) error
	GetTitsWithReportsThreshold(ctx context.Context, reportsThreshold int) ([]domain.Tits, error)
}

type Storage interface {
	CreateImageFromFile(ctx context.Context, imageName string, filePath string) error
	CreateImageFromBytes(ctx context.Context, imageName string, imageData []byte) error
	GetImageUrl(imageID string) string
}
