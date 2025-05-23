package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/repository"
	"time"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) repository.TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, userID int64, req model.CreateTaskRequest) (*model.Task, error) {
	status := req.Status
	if status == "" {
		status = model.TaskStatusPending
	}

	now := time.Now()
	task := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		UserID:      userID,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO tasks (title, description, status, user_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания задачи: %w", err)
	}

	return &task, nil
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
