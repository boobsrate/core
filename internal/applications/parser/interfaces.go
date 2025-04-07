package parser

import (
	"context"

	"github.com/boobsrate/core/internal/domain"
)

type DetectorService interface {
	Detect(ctx context.Context, url string) (domain.DetectionResult, error)
}

type TitsService interface {
	CreateTitsFromFile(ctx context.Context, filename, filePath string) error
	CreateTitsFromBytes(ctx context.Context, filename string, file []byte, url string) error
}

type TaskRepo interface {
	GetTask(ctx context.Context) (domain.Task, error)
	CreateTask(ctx context.Context, task []domain.Task) error
	UpdateTask(ctx context.Context, task domain.Task) error
	GetCountUnprocessedTasks(ctx context.Context) (int, error)
}
