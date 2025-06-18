package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kkboranbay/task-service/internal/mocks"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/service"
	"github.com/kkboranbay/task-service/internal/testutils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TaskHandlerTestSuite struct {
	suite.Suite
	handler     *TaskHandler
	mockRepo    *mocks.MockTaskRepository
	taskService *service.TaskService
	router      *gin.Engine
}

func (suite *TaskHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *TaskHandlerTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockTaskRepository)
	logger := zerolog.Nop()
	suite.taskService = service.NewTaskService(suite.mockRepo, &logger)
	suite.handler = NewTaskHandler(suite.taskService, &logger)

	suite.router = gin.New()

	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})

	api := suite.router.Group("/api/v1")
	suite.handler.Register(api)
}

func (suite *TaskHandlerTestSuite) TestCreateTask() {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful_creation",
			requestBody: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = "task title"
				r.Description = "task description"
				r.Status = "pending"
			}),
			setupMock: func() {
				expectedTask := testutils.TaskFixture(func(t *model.Task) {
					t.Title = "task title"
					t.Description = "task description"
					t.Status = "pending"
				})
				suite.mockRepo.On("Create", mock.Anything, int64(1), mock.AnythingOfType("model.CreateTaskRequest")).
					Return(expectedTask, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"title":       "task title",
				"description": "task description",
				"status":      "pending",
			},
		},
		{
			name:           "invalid_json",
			requestBody:    "invalid json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(400),
				"message": "некорректные данные запроса",
			},
		},
		{
			name: "empty_title",
			requestBody: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = ""
				r.Description = "task description"
				r.Status = "pending"
			}),
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(400),
				"message": "некорректные данные запроса",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			var bodyReader *bytes.Reader
			if str, ok := tt.requestBody.(string); ok {
				bodyReader = bytes.NewReader([]byte(str))
			} else {
				bodyBytes, _ := json.Marshal(tt.requestBody)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bodyReader)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(suite.T(), err)

			for key, value := range tt.expectedBody {
				assert.Equal(suite.T(), value, response[key], "Field %s mismatch", key)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *TaskHandlerTestSuite) TestGetTask() {
	tests := []struct {
		name           string
		taskID         string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful_get",
			taskID: "1",
			setupMock: func() {
				expectedTask := testutils.TaskFixture()
				suite.mockRepo.On("GetByID", mock.Anything, int64(1), int64(1)).
					Return(expectedTask, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"title":       "Test task",
				"description": "Test Description",
				"status":      "pending",
			},
		},
		{
			name:           "invalid_id",
			taskID:         "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(400),
				"message": "некорректный ID задачи",
			},
		},
		{
			name:   "task_not_found",
			taskID: "999",
			setupMock: func() {
				suite.mockRepo.On("GetByID", mock.Anything, int64(999), int64(1)).
					Return(nil, fmt.Errorf("task not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"code":    float64(404),
				"message": "задача не найдена",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(suite.T(), err)

			for key, value := range tt.expectedBody {
				assert.Equal(suite.T(), value, response[key], "Field %s mismatch", key)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func TestTaskHandlerSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlerTestSuite))
}
