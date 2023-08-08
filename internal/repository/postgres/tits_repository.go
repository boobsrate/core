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
	titsModels, err := t.createTitsSlice()
	if err != nil {
		return nil, err
	}
	err = t.executeQuery(ctx, &titsModels)
	if err != nil {
		return nil, err
	}
	tits := t.convertToDomain(titsModels)
	return tits, nil
}

func (t *TitsRepository) createTitsSlice() ([]titsModel, error) {
	return make([]titsModel, 0, 2), nil
}

func (t *TitsRepository) executeQuery(ctx context.Context, titsModels *[]titsModel) error {
	return t.db.NewSelect().
		Model(titsModels).
		OrderExpr("random()").
		Limit(2).
		Scan(ctx)
}

func (t *TitsRepository) convertToDomain(titsModels []titsModel) []domain.Tits {
	return titsModelsToDomain(titsModels)
}

func (t *TitsRepository) CreateTits(ctx context.Context, tits domain.Tits) error {
	model, err := t.createModel(tits)
	if err != nil {
		return err
	}
	err = t.executeInsert(ctx, model)
	return err
}

func (t *TitsRepository) createModel(tits domain.Tits) (*titsModel, error) {
	model := &titsModel{}
	model.FromDomain(tits)
	return model, nil
}

func (t *TitsRepository) executeInsert(ctx context.Context, model *titsModel) error {
	_, err := t.db.NewInsert().
		Model(model).
		Exec(ctx)
	return err
}

func (t *TitsRepository) IncreaseRating(ctx context.Context, titsID string) (int64, error) {
	rating, err := t.executeUpdate(ctx, titsID)
	if err != nil {
		return rating, err
	}
	return rating, nil
}

func (t *TitsRepository) executeUpdate(ctx context.Context, titsID string) (int64, error) {
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
	return rating, err
}