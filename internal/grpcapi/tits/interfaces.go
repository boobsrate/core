package titspbv1

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type Service interface {
	GetTits(ctx context.Context) ([]domain.Tits, error)
	IncreaseRating(ctx context.Context, titsID string) error
}
