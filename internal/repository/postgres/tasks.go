package postgres

import (
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/uptrace/bun"
)

type tasksModel struct {
	bun.BaseModel `bun:"table:tasks,alias:tasks,select:tasks"`

	ID              string                 `bun:"id,pk"`
	CreatedAt       time.Time              `bun:"created_at"`
	Processed       bool                   `bun:"processed"`
	Url             string                 `bun:"url"`
	Status          string                 `bun:"status"`
	NeedRetry       bool                   `bun:"need_retry"`
	Error           string                 `bun:"error"`
	DetectionResult domain.DetectionResult `bun:"detection_result,type:jsonb"`
}

func (t *tasksModel) FromDomain(task domain.Task) {
	t.ID = task.ID
	t.Processed = task.Processed
	t.CreatedAt = task.CreatedAt
	t.Url = task.Url
	t.Status = task.Status
	t.NeedRetry = task.NeedRetry
	t.Error = task.Error
	t.DetectionResult = task.DetectionResult
}

func TaskModelToDomain(model tasksModel) domain.Task {
	return domain.Task{
		ID:              model.ID,
		CreatedAt:       model.CreatedAt,
		Processed:       model.Processed,
		Url:             model.Url,
		Status:          model.Status,
		NeedRetry:       model.NeedRetry,
		Error:           model.Error,
		DetectionResult: model.DetectionResult,
	}
}

func tasksModelsToDomain(models []tasksModel) []domain.Task {
	tasks := make([]domain.Task, 0, len(models))
	for _, model := range models {
		tasks = append(tasks, TaskModelToDomain(model))
	}
	return tasks
}
