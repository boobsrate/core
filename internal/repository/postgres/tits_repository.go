package postgres

import (
	"context"
	"database/sql"
	"time"

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

func (t *TitsRepository) GetTop(ctx context.Context, limit int, abyss bool) ([]domain.Tits, error) {
	titsModels := make([]titsModel, 0, limit)
	err := t.db.NewSelect().
		Model(&titsModels).
		Where("COALESCE(abyss, FALSE) = ?", abyss).
		OrderExpr("rating DESC").
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	tits := titsModelsToDomain(titsModels)
	return tits, nil
}

func (t *TitsRepository) GetTits(ctx context.Context) ([]domain.Tits, error) {
	titsModels := make([]titsModel, 0, 2)
	err := t.db.NewSelect().
		Model(&titsModels).
		Where("COALESCE(abyss, FALSE) = ?", false).
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
		if err != nil {
			return err
		}
		_, err = tx.NewInsert().
			Model(&voteModel{TitsID: titsID, CreatedAt: time.Now().UTC()}).
			Exec(ctx)
		return err
	})
	if err != nil {
		return rating, err
	}
	return rating, nil
}

func (t *TitsRepository) Report(ctx context.Context, titsID string) error {
	_, err := t.db.NewInsert().
		Model(&reportModel{TitsID: titsID, CreatedAt: time.Now().UTC()}).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TitsRepository) GetReportsCount(ctx context.Context, titsID string) (int, error) {
	count, err := t.db.NewSelect().
		Model(&reportModel{}).
		Where("tits_id = ?", titsID).
		Count(ctx)
	return count, err
}

func (t *TitsRepository) GetTitsWithReportsThreshold(ctx context.Context, reportsThreshold int) ([]domain.Tits, error) {
	var titsModels []titsModel
	err := t.db.NewSelect().
		Model(&titsModels).
		Where("COALESCE(abyss, FALSE) = ?", false).
		Join("JOIN reports ON reports.tits_id = tits.id").
		Group("tits.id").
		Having("COUNT(tits.id) >= ?", reportsThreshold).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	tits := titsModelsToDomain(titsModels)
	return tits, nil
}

func (t *TitsRepository) MoveToAbyss(ctx context.Context, titsID string) error {
	_, err := t.db.NewUpdate().
		Model(&titsModel{}).
		Set("abyss = ?", true).
		Where("id = ?", titsID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
