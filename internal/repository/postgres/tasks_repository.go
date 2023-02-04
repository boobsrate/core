package postgres

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type TasksRepository struct {
	db *bun.DB
}

func NewTasksRepository(db *bun.DB) *TasksRepository {
	return &TasksRepository{
		db: db,
	}
}

func (r *TasksRepository) GetTask(ctx context.Context) (domain.Task, error) {
	var task tasksModel
	err := r.db.NewSelect().Model(&task).Where("processed = False").Limit(1).Scan(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	return TaskModelToDomain(task), nil
}

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) error {
	model := tasksModel{}
	model.FromDomain(task)
	_, err := r.db.NewInsert().Model(&model).Exec(ctx)
	return err
}

func (r *TasksRepository) UpdateTask(ctx context.Context, task domain.Task) error {
	model := tasksModel{}
	model.FromDomain(task)
	_, err := r.db.NewUpdate().Model(&model).WherePK().Exec(ctx)
	return err
}
