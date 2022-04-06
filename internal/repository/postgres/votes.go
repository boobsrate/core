package postgres

import (
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type voteModel struct {
	bun.BaseModel `bun:"table:votes"`

	TitsID    string    `bun:"tits_id"`
	CreatedAt time.Time `bun:"created_at"`
}

func (v *voteModel) FromDomain(vote domain.Vote) {
	v.CreatedAt = vote.CreatedAt
	v.TitsID = vote.TitsID
}

func voteModelToDomain(model voteModel) domain.Vote {
	return domain.Vote{
		TitsID:    model.TitsID,
		CreatedAt: model.CreatedAt,
	}
}
