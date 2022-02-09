package postgres

import (
	"context"
	"database/sql"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type TitsRepository struct {
	db *bun.DB
}

func NewTitsRepository(db *bun.DB) *TitsRepository {
	return &TitsRepository{
		db: db,
	}
}

func (t *TitsRepository) GetTits(ctx context.Context) ([]domain.Tits, error) {
	titsModels := make([]titsModel, 0, 2)
	err := t.db.NewSelect().
		Model(&titsModels).
		OrderExpr("random()").
		Limit(2).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	tits := titsModelsToDomain(titsModels)
	return tits, nil
}

func (t *TitsRepository) CreateTits(ctx context.Context, tits domain.Tits) error {
	model := titsModel{}
	model.FromDomain(tits)
	_, err := t.db.NewInsert().
		Model(&model).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TitsRepository) IncreaseRating(ctx context.Context, titsID string) (int64, error) {
	var rating int64
	err := t.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().
			Model(&titsModel{}).
			Set("rating = rating + 1").
			Where("id = ?", titsID).
			Returning("rating").
			Exec(ctx, &rating)
		return err
	})
	if err != nil {
		return rating, err
	}
	return rating, nil
}
