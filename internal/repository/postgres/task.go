package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	query := `
		SELECT id, title, description, status, user_id, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	var task model.Task

	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) List(ctx context.Context, userID int64, limit, offset int) (*model.TaskListResponse, error) {
	countQuery := `SELECT count(*) FROM tasks WHERE user_id = $1`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("ошибка подсчета задач: %w", err)
	}

	query := `
		SELECT id, title, description, status, user_id, due_date, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %w", err)
	}

	return &model.TaskListResponse{
		Total: total,
		Tasks: tasks,
	}, nil
}

func (r *TaskRepository) Update(ctx context.Context, id, userID int64, req model.UpdateTaskRequest) (*model.Task, error) {
	task, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
		RETURNING id
	`

	err = r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
		task.UserID,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
