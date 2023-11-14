package abyss

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type Service interface {
	MoveToAbyss(ctx context.Context, titsID string) error
	GetTitsWithReportsThreshold(ctx context.Context, reportsThreshold int) ([]domain.Tits, error)
}
