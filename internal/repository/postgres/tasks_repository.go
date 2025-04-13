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
	err := r.db.NewSelect().Model(&task).Where("processed = False").OrderExpr("random()").Limit(1).Scan(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	return TaskModelToDomain(task), nil
}

func (r *TasksRepository) GetCountUnprocessedTasks(ctx context.Context) (int, error) {
	var count int
	count, err := r.db.NewSelect().Model((*tasksModel)(nil)).Where("processed = False").Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TasksRepository) CreateTask(ctx context.Context, task []domain.Task) error {
	var models []tasksModel
	models = make([]tasksModel, len(task))
	for i, t := range task {
		models[i].FromDomain(t)
	}
	_, err := r.db.NewInsert().Model(&models).On("CONFLICT (url) DO nothing").Exec(ctx)
	return err
}

func (r *TasksRepository) UpdateTask(ctx context.Context, task domain.Task) error {
	model := tasksModel{}
	model.FromDomain(task)
	_, err := r.db.NewUpdate().Model(&model).WherePK().Exec(ctx)
	return err
}
