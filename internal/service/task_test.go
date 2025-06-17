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

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceTestSuite))
}
