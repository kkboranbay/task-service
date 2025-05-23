package repository

import (
	"context"
	"github.com/kkboranbay/task-service/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, userID int64, task model.CreateTaskRequest) (*model.Task, error)
	GetByID(ctx context.Context, id, userID int64) (*model.Task, error)
	List(ctx context.Context, userID int64, limit, offset int) (*model.TaskListResponse, error)
	Update(ctx context.Context, id, userID int64, task model.UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id, userID int64) error
}

type Repository struct {
	Task TaskRepository
}
