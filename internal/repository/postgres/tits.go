package postgres

import (
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type titsModel struct {
	bun.BaseModel `bun:"table:tits,alias:tits,select:tits"`

	ID        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:"created_at"`
	Rating    int64     `bun:"rating"`
	Abyss     bool      `bun:"abyss"`
}

func (t *titsModel) FromDomain(tits domain.Tits) {
	t.CreatedAt = tits.CreatedAt
	t.Rating = tits.Rating
	t.ID = tits.ID
	t.Abyss = tits.Abyss
}

func titsModelToDomain(model titsModel) domain.Tits {
	return domain.Tits{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		Rating:    model.Rating,
		Abyss:     model.Abyss,
	}
}

func titsModelsToDomain(models []titsModel) []domain.Tits {
	tits := make([]domain.Tits, 0, len(models))
	for _, model := range models {
		tits = append(tits, titsModelToDomain(model))
	}
	return tits
}
