package initiator

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type TitsService interface {
	CreateTitsFromFile(ctx context.Context, filename, filePath string) error
	CreateTitsFromBytes(ctx context.Context, filename string, file []byte) error
}

type TaskRepo interface {
	GetTask(ctx context.Context) (domain.Task, error)
	CreateTask(ctx context.Context, task []domain.Task) error
	UpdateTask(ctx context.Context, task domain.Task) error
	GetCountUnprocessedTasks(ctx context.Context) (int, error)
}
