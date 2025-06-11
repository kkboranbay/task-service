package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/repository"
	"github.com/rs/zerolog"
)

type TaskService struct {
	repo repository.TaskRepository
	log  *zerolog.Logger
}

func NewTaskService(repo repository.TaskRepository, log *zerolog.Logger) *TaskService {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID int64, req model.CreateTaskRequest) (*model.Task, error) {
	if req.Title == "" {
		return nil, errors.New("отсутствует заголовок задачи")
	}

	s.log.Info().Int64("user_id", userID).Str("title", req.Title).Msg("создание новой задачи")

	task, err := s.repo.Create(ctx, userID, req)
	if err != nil {
		s.log.Error().Err(err).Int64("user_id", userID).Str("title", req.Title).Msg("ошибка создания задачи")
		return nil, fmt.Errorf("не удалось создать задачу: %w", err)
	}

	s.log.Info().Int64("task_id", task.ID).Int64("user_id", userID).Msg("задача успешно создана")
	return task, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id, userID int64) (*model.Task, error) {
	s.log.Info().Int64("task_id", id).Int64("user_id", userID).Msg("получение задачи по ID")

	task, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		s.log.Error().Err(err).Int64("task_id", id).Int64("user_id", userID).Msg("ошибка получения задачи")
		return nil, fmt.Errorf("не удалось получить задачу: %w", err)
	}

	return task, nil
}

func (s *TaskService) GetTaskList(ctx context.Context, userID int64, page, pageSize int) (*model.TaskListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	s.log.Info().Int64("user_id", userID).Int("page", page).Int("page_size", pageSize).Msg("получение списка задач")

	resp, err := s.repo.List(ctx, userID, pageSize, offset)
	if err != nil {
		s.log.Error().Err(err).Int64("user_id", userID).Msg("ошибка получения списка задач")
		return nil, fmt.Errorf("не удалось получить список задач: %w", err)
	}

	return resp, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id, userID int64, req model.UpdateTaskRequest) (*model.Task, error) {
	s.log.Info().Int64("task_id", id).Int64("user_id", userID).Msg("обновление задачи")

	if req.Status != nil {
		status := *req.Status
		if status != model.TaskStatusPending && status != model.TaskStatusInProgress && status != model.TaskStatusCompleted {
			return nil, errors.New("некорректный статус задачи")
		}
	}

	task, err := s.repo.Update(ctx, id, userID, req)
	if err != nil {
		s.log.Error().Err(err).Int64("task_id", id).Int64("user_id", userID).Msg("ошибка обновления задачи")
		return nil, fmt.Errorf("не удалось обновить задачу: %w", err)
	}

	s.log.Info().Int64("task_id", id).Int64("user_id", userID).Msg("задача успешно обновлена")
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id, userID int64) error {
	s.log.Info().Int64("task_id", id).Int64("user_id", userID).Msg("удаление задачи")

	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		s.log.Error().Err(err).Int64("task_id", id).Int64("user_id", userID).Msg("ошибка удаления задачи")
		return fmt.Errorf("не удалось удалить задачу: %w", err)
	}

	s.log.Info().Int64("task_id", id).Int64("user_id", userID).Msg("задача успешно удалена")
	return nil
}
