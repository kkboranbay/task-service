package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kkboranbay/task-service/internal/api/handler"
	"github.com/kkboranbay/task-service/internal/api/middleware"
	"github.com/kkboranbay/task-service/internal/config"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/repository/postgres"
	"github.com/kkboranbay/task-service/internal/service"
	"github.com/kkboranbay/task-service/internal/testutils"
	"github.com/kkboranbay/task-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type E2ETestSuite struct {
	suite.Suite
	testDB     *testutils.TestDB
	server     *httptest.Server
	jwtToken   string
	httpClient *http.Client
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.testDB = testutils.NewTestDB(suite.T())

	cfg := config.Config{
		Server: config.ServerConfig{
			Port:            "8080",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-key",
			TokenExpireDelta: 24 * time.Hour,
		},
		Logger: config.LoggerConfig{
			Level: "error",
		},
	}

	logger.SetupLogger(cfg.Logger)
	log := logger.L()

	taskRepo := postgres.NewTaskRepository(suite.testDB.Pool)
	taskService := service.NewTaskService(taskRepo, log)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	jwtMiddleware := middleware.NewJWTMiddleware(cfg.Auth, log)
	authHandler := handler.NewAuthHandler(jwtMiddleware, log)
	authHandler.Register(router)

	apiGroup := router.Group("/api/v1")
	apiGroup.Use(jwtMiddleware.AuthRequired())

	taskHandler := handler.NewTaskHandler(taskService, log)
	taskHandler.Register(apiGroup)

	suite.server = httptest.NewServer(router)
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.jwtToken = suite.getJWTToken()
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.testDB != nil {
		suite.testDB.Close(suite.T())
	}
}

func (suite *E2ETestSuite) SetupTest() {
	suite.testDB.Truncate(suite.T())
}

func (suite *E2ETestSuite) getJWTToken() string {
	loginReq := testutils.LoginRequestFixture()
	loginData, _ := json.Marshal(loginReq)

	resp, err := suite.httpClient.Post(
		suite.server.URL+"/auth/login",
		"application/json",
		bytes.NewReader(loginData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var loginResp model.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(suite.T(), err)

	return loginResp.Token
}

func (suite *E2ETestSuite) makeAuthenticatedRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, suite.server.URL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+suite.jwtToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return suite.httpClient.Do(req)
}

func (suite *E2ETestSuite) TestTaskCRUDFlow() {
	// 1. Создаем задачу
	createReq := testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
		r.Title = "E2E Test Task"
		r.Description = "testing"
		r.Status = model.TaskStatusPending
	})

	resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/tasks", createReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdTask model.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	require.NoError(suite.T(), err)

	assert.Greater(suite.T(), createdTask.ID, int64(0))
	assert.Equal(suite.T(), createReq.Title, createdTask.Title)
	assert.Equal(suite.T(), createReq.Description, createdTask.Description)
	assert.Equal(suite.T(), model.TaskStatusPending, createdTask.Status)

	// 2. Получаем созданную задачу
	resp, err = suite.makeAuthenticatedRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", createdTask.ID), nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedTask model.Task
	err = json.NewDecoder(resp.Body).Decode(&retrievedTask)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), createdTask.ID, retrievedTask.ID)
	assert.Equal(suite.T(), createdTask.Title, retrievedTask.Title)
	assert.Equal(suite.T(), createdTask.Description, retrievedTask.Description)

	// 3. Обновляем задачу
	updateReq := testutils.UpdateTaskRequestFixture(func(r *model.UpdateTaskRequest) {
		r.Title = testutils.StringPtr("Updated E2E Test Task")
		r.Status = testutils.TaskStatusPtr(model.TaskStatusInProgress)
	})

	resp, err = suite.makeAuthenticatedRequest("PUT", fmt.Sprintf("/api/v1/tasks/%d", createdTask.ID), updateReq)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var updatedTask model.Task
	err = json.NewDecoder(resp.Body).Decode(&updatedTask)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), createdTask.ID, updatedTask.ID)
	assert.Equal(suite.T(), *updateReq.Title, updatedTask.Title)
	assert.Equal(suite.T(), *updateReq.Status, updatedTask.Status)

	// 4. Получаем список задач
	resp, err = suite.makeAuthenticatedRequest("GET", "/api/v1/tasks", nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var taskList model.TaskListResponse
	err = json.NewDecoder(resp.Body).Decode(&taskList)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int64(1), taskList.Total)
	assert.Len(suite.T(), taskList.Tasks, 1)
	assert.Equal(suite.T(), updatedTask.ID, taskList.Tasks[0].ID)

	// 5. Удаляем задачу
	resp, err = suite.makeAuthenticatedRequest("DELETE", fmt.Sprintf("/api/v1/tasks/%d", createdTask.ID), nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// 6. Проверяем, что задача удалена
	resp, err = suite.makeAuthenticatedRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", createdTask.ID), nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *E2ETestSuite) TestTaskListPagination() {
	const taskCount = 15
	createdTasks := make([]*model.Task, taskCount)

	for i := 0; i < taskCount; i++ {
		createReq := testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
			r.Title = fmt.Sprintf("Task %d", i+1)
			r.Description = fmt.Sprintf("Description for task %d", i+1)
			r.Status = model.TaskStatusPending
		})

		resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/tasks", createReq)
		require.NoError(suite.T(), err)
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var task model.Task
		err = json.NewDecoder(resp.Body).Decode(&task)
		require.NoError(suite.T(), err)
		createdTasks[i] = &task
	}

	resp, err := suite.makeAuthenticatedRequest("GET", "/api/v1/tasks?page=1&page_size=10", nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var page1 model.TaskListResponse
	err = json.NewDecoder(resp.Body).Decode(&page1)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int64(taskCount), page1.Total)
	assert.Len(suite.T(), page1.Tasks, 10)

	resp, err = suite.makeAuthenticatedRequest("GET", "/api/v1/tasks?page=2&page_size=10", nil)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var page2 model.TaskListResponse
	err = json.NewDecoder(resp.Body).Decode(&page2)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), int64(taskCount), page2.Total)
	assert.Len(suite.T(), page2.Tasks, 5)
}

func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	createReq := testutils.CreateTaskRequestFixture()
	reqBody, _ := json.Marshal(createReq)

	resp, err := suite.httpClient.Post(
		suite.server.URL+"/api/v1/tasks",
		"application/json",
		bytes.NewReader(reqBody),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *E2ETestSuite) TestInvalidToken() {
	req, err := http.NewRequest("GET", suite.server.URL+"/api/v1/tasks", nil)
	require.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *E2ETestSuite) TestValidationErrors() {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "empty_title",
			requestBody: map[string]interface{}{
				"title":       "",
				"description": "Valid description",
				"status":      model.TaskStatusPending,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing_title",
			requestBody: map[string]interface{}{
				"description": "Valid description",
				"status":      model.TaskStatusPending,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty_body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/tasks", tt.requestBody)
			require.NoError(suite.T(), err)
			defer resp.Body.Close()

			assert.Equal(suite.T(), tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
