package service

import (
	"context"
	"errors"
	"github.com/kkboranbay/task-service/internal/mocks"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/testutils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TaskServiceTestSuite struct {
	suite.Suite
	service  *TaskService
	mockRepo *mocks.MockTaskRepository
	ctx      context.Context
}

func (suite *TaskServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockTaskRepository)
	logger := zerolog.Nop()
	suite.service = NewTaskService(suite.mockRepo, &logger)
	suite.ctx = context.Background()
}

func (suite *TaskServiceTestSuite) TestCreateTask() {
	tests := []struct {
		name       string
		userID     int64
		req        model.CreateTaskRequest
		setupMock  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:   "successful_creation",
			userID: 1,
			req:    testutils.CreateTaskRequestFixture(),
			setupMock: func() {
				expectedTask := testutils.TaskFixture()
				suite.mockRepo.On("Create", suite.ctx, int64(1), mock.AnythingOfType("model.CreateTaskRequest")).
					Return(expectedTask, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "empty_title",
			userID: 1,
			req: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = ""
			}),
			setupMock:  func() {},
			wantErr:    true,
			wantErrMsg: "отсутствует заголовок задачи",
		},
		{
			name:   "repository_error",
			userID: 1,
			req:    testutils.CreateTaskRequestFixture(),
			setupMock: func() {
				suite.mockRepo.On("Create", suite.ctx, int64(1), mock.AnythingOfType("model.CreateTaskRequest")).
					Return(nil, errors.New("database error")).Once()
			},
			wantErr:    true,
			wantErrMsg: "не удалось создать задачу",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			task, err := suite.service.CreateTask(suite.ctx, tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tt.wantErrMsg)
				assert.Nil(suite.T(), task)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), task)
				assert.Equal(suite.T(), tt.userID, task.UserID)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *TaskServiceTestSuite) TestGetTaskByID() {
	tests := []struct {
		name       string
		id         int64
		userID     int64
		setupMock  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:   "successful_get",
			id:     1,
			userID: 1,
			setupMock: func() {
				expectedTask := testutils.TaskFixture()
				suite.mockRepo.On("GetByID", suite.ctx, int64(1), int64(1)).
					Return(expectedTask, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "task_not_found",
			id:     999,
			userID: 1,
			setupMock: func() {
				suite.mockRepo.On("GetByID", suite.ctx, int64(999), int64(1)).
					Return(nil, errors.New("task not found")).Once()
			},
			wantErr:    true,
			wantErrMsg: "не удалось получить задачу",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			task, err := suite.service.GetTaskByID(suite.ctx, tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tt.wantErrMsg)
				assert.Nil(suite.T(), task)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), task)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *TaskServiceTestSuite) TestGetTaskList() {
	tests := []struct {
		name      string
		userID    int64
		page      int
		pageSize  int
		setupMock func()
		wantErr   bool
	}{
		{
			name:     "successful_get_first_page",
			userID:   1,
			page:     1,
			pageSize: 10,
			setupMock: func() {
				expectedResponse := &model.TaskListResponse{
					Total: 1,
					Tasks: []model.Task{*testutils.TaskFixture()},
				}
				suite.mockRepo.On("List", suite.ctx, int64(1), 10, 0).
					Return(expectedResponse, nil).Once()
			},
		},
		{
			name:     "successful_get_second_page",
			userID:   1,
			page:     2,
			pageSize: 5,
			setupMock: func() {
				expectedResponse := &model.TaskListResponse{
					Tasks: []model.Task{*testutils.TaskFixture()},
					Total: 10,
				}
				suite.mockRepo.On("List", suite.ctx, int64(1), 5, 5).
					Return(expectedResponse, nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "invalid_page_defaults_to_1",
			userID:   1,
			page:     0,
			pageSize: 10,
			setupMock: func() {
				expectedResponse := &model.TaskListResponse{Tasks: []model.Task{}, Total: 0}
				suite.mockRepo.On("List", suite.ctx, int64(1), 10, 0).
					Return(expectedResponse, nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "invalid_page_size_defaults_to_10",
			userID:   1,
			page:     1,
			pageSize: 0,
			setupMock: func() {
				expectedResponse := &model.TaskListResponse{Tasks: []model.Task{}, Total: 0}
				suite.mockRepo.On("List", suite.ctx, int64(1), 10, 0).
					Return(expectedResponse, nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "page_size_too_large_capped_at_100",
			userID:   1,
			page:     1,
			pageSize: 200,
			setupMock: func() {
				expectedResponse := &model.TaskListResponse{Tasks: []model.Task{}}
				suite.mockRepo.On("List", suite.ctx, int64(1), 100, 0).
					Return(expectedResponse, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			result, err := suite.service.GetTaskList(suite.ctx, tt.userID, tt.page, tt.pageSize)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), result)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *TaskServiceTestSuite) TestUpdateTask() {
	tests := []struct {
		name       string
		id         int64
		userID     int64
		req        model.UpdateTaskRequest
		setupMock  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:   "successful_update",
			id:     1,
			userID: 1,
			req:    testutils.UpdateTaskRequestFixture(),
			setupMock: func() {
				expectedTask := testutils.TaskFixture()
				suite.mockRepo.On("Update", suite.ctx, int64(1), int64(1), mock.AnythingOfType("model.UpdateTaskRequest")).
					Return(expectedTask, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "invalid_status",
			id:     1,
			userID: 1,
			req: testutils.UpdateTaskRequestFixture(func(r *model.UpdateTaskRequest) {
				invalidStatus := model.TaskStatus("invalid")
				r.Status = &invalidStatus
			}),
			setupMock:  func() {},
			wantErr:    true,
			wantErrMsg: "некорректный статус задачи",
		},
		{
			name:   "repository_error",
			id:     1,
			userID: 1,
			req:    testutils.UpdateTaskRequestFixture(),
			setupMock: func() {
				suite.mockRepo.On("Update", suite.ctx, int64(1), int64(1), mock.AnythingOfType("model.UpdateTaskRequest")).
					Return(nil, errors.New("database error")).Once()
			},
			wantErr:    true,
			wantErrMsg: "не удалось обновить задачу",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			task, err := suite.service.UpdateTask(suite.ctx, tt.id, tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tt.wantErrMsg)
				assert.Nil(suite.T(), task)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), task)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceTestSuite))
}
