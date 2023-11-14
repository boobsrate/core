package tits

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type Service interface {
	GetTits(ctx context.Context) ([]domain.Tits, error)
	GetTop(ctx context.Context, limit int, abyss bool) ([]domain.Tits, error)
	IncreaseRating(ctx context.Context, titsID string) error
	Report(ctx context.Context, titsID string) error
}
