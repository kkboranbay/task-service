package mocks

import (
	"context"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, userID int64, task model.CreateTaskRequest) (*model.Task, error) {
	args := m.Called(ctx, userID, task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id, userID int64) (*model.Task, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) List(ctx context.Context, userID int64, limit, offset int) (*model.TaskListResponse, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TaskListResponse), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, id, userID int64, req model.UpdateTaskRequest) (*model.Task, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
