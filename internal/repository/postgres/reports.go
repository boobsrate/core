package postgres

import (
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type reportModel struct {
	bun.BaseModel `bun:"table:reports"`

	TitsID    string    `bun:"tits_id"`
	CreatedAt time.Time `bun:"created_at"`
}

func (v *reportModel) FromDomain(report domain.Report) {
	v.CreatedAt = report.CreatedAt
	v.TitsID = report.TitsID
}

func reportModelToDomain(model reportModel) domain.Report {
	return domain.Report{
		TitsID:    model.TitsID,
		CreatedAt: model.CreatedAt,
	}
}
