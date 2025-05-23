package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/repository"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) repository.TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, userID int64, task model.CreateTaskRequest) (*model.Task, error) {
	panic("implement me")
}

func (r *TaskRepository) GetByID(ctx context.Context, id, userID int64) (*model.Task, error) {
	panic("implement me")
}

func (r *TaskRepository) List(ctx context.Context, userID int64, limit, offset int) (*model.TaskListResponse, error) {
	panic("implement me")
}

func (r *TaskRepository) Update(ctx context.Context, id, userID int64, task model.UpdateTaskRequest) (*model.Task, error) {
	panic("implement me")
}

func (r *TaskRepository) Delete(ctx context.Context, id, userID int64) error {
	panic("implement me")
}
