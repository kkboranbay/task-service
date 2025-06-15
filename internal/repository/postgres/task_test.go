package postgres

import (
	"context"
	"fmt"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
	testDB *testutils.TestDB
	repo   *TaskRepository
	ctx    context.Context
}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *TaskRepositoryTestSuite) SetupSuite() {
	suite.testDB = testutils.NewTestDB(suite.T())
	suite.repo = &TaskRepository{pool: suite.testDB.Pool}
	suite.ctx = context.Background()
}

// TearDownSuite выполняется один раз после всех тестов
func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	suite.testDB.Close(suite.T())
}

// SetupTest выполняется перед каждым тестом
func (suite *TaskRepositoryTestSuite) SetupTest() {
	suite.testDB.Truncate(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		userID  int64
		req     model.CreateTaskRequest
		wantErr bool
	}{
		{
			name:   "successful_creation",
			userID: 1,
			req: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = "Test Task"
				r.Description = "Test Description"
			}),
			wantErr: false,
		},
		{
			name:   "empty_title",
			userID: 1,
			req: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = ""
				r.Description = "Test Description"
			}),
			wantErr: true,
		},
		{
			name:   "with_empty_description",
			userID: 1,
			req: testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
				r.Title = "Test Task"
				r.Description = ""
			}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			task, err := suite.repo.Create(suite.ctx, tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), task)
				return
			}

			assert.Greater(suite.T(), task.ID, int64(0))
			assert.Equal(suite.T(), tt.userID, task.UserID)
			assert.Equal(suite.T(), tt.req.Title, task.Title)
			assert.Equal(suite.T(), tt.req.Description, task.Description)
			assert.Equal(suite.T(), model.TaskStatusPending, task.Status)
			assert.WithinDuration(suite.T(), time.Now(), task.CreatedAt, 5*time.Second)
			assert.WithinDuration(suite.T(), time.Now(), task.UpdatedAt, 5*time.Second)
		})
	}
}

func (suite *TaskRepositoryTestSuite) TestGetByID() {
	userID := int64(1)
	req := testutils.CreateTaskRequestFixture()
	createdTask, err := suite.repo.Create(suite.ctx, userID, req)
	require.NoError(suite.T(), err)

	tests := []struct {
		name    string
		id      int64
		userID  int64
		wantErr bool
	}{
		{
			name:    "existing_task",
			id:      createdTask.ID,
			userID:  userID,
			wantErr: false,
		},
		{
			name:    "non_existing_task",
			id:      999,
			userID:  userID,
			wantErr: true,
		},
		{
			name:    "wrong_user",
			id:      createdTask.ID,
			userID:  999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			task, err := suite.repo.GetByID(suite.ctx, tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), task)
				return
			}

			require.NoError(suite.T(), err)
			require.NotNil(suite.T(), task)

			assert.Equal(suite.T(), tt.id, task.ID)
			assert.Equal(suite.T(), tt.userID, task.UserID)
			assert.Equal(suite.T(), createdTask.Title, task.Title)
			assert.Equal(suite.T(), createdTask.Description, task.Description)
			assert.Equal(suite.T(), createdTask.Status, task.Status)
		})
	}
}

func (suite *TaskRepositoryTestSuite) TestList() {
	userID := int64(1)
	otherUserID := int64(2)

	for i := 0; i < 5; i++ {
		req := testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
			r.Title = fmt.Sprintf("Task %d", i+1)
		})
		_, err := suite.repo.Create(suite.ctx, userID, req)
		require.NoError(suite.T(), err)
	}

	req := testutils.CreateTaskRequestFixture(func(r *model.CreateTaskRequest) {
		r.Title = "Other User Task"
	})
	_, err := suite.repo.Create(suite.ctx, otherUserID, req)
	require.NoError(suite.T(), err)

	tests := []struct {
		name      string
		userID    int64
		limit     int
		offset    int
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "get_all_tasks",
			userID:    userID,
			limit:     10,
			offset:    0,
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "pagination_first_page",
			userID:    userID,
			limit:     3,
			offset:    0,
			wantCount: 3,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "pagination_second_page",
			userID:    userID,
			limit:     3,
			offset:    3,
			wantCount: 2,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "other_user_tasks",
			userID:    otherUserID,
			limit:     10,
			offset:    0,
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "non_existing_user",
			userID:    999,
			limit:     10,
			offset:    0,
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := suite.repo.List(suite.ctx, tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				return
			}

			require.NoError(suite.T(), err)
			require.NotNil(suite.T(), result)

			assert.Len(suite.T(), result.Tasks, tt.wantCount)
			assert.Equal(suite.T(), tt.wantTotal, result.Total)

			for _, task := range result.Tasks {
				assert.Equal(suite.T(), tt.userID, task.UserID)
			}
		})
	}
}

func (suite *TaskRepositoryTestSuite) TestUpdate() {
	userID := int64(1)
	req := testutils.CreateTaskRequestFixture()
	createdTask, err := suite.repo.Create(suite.ctx, userID, req)
	require.NoError(suite.T(), err)

	tests := []struct {
		name    string
		id      int64
		userID  int64
		req     model.UpdateTaskRequest
		wantErr bool
	}{
		{
			name:   "update_title",
			id:     createdTask.ID,
			userID: userID,
			req: testutils.UpdateTaskRequestFixture(func(r *model.UpdateTaskRequest) {
				r.Title = testutils.StringPtr("Updated Title")
				r.Description = nil
				r.Status = nil
			}),
			wantErr: false,
		},
		{
			name:   "update_all_fields",
			id:     createdTask.ID,
			userID: userID,
			req: testutils.UpdateTaskRequestFixture(func(r *model.UpdateTaskRequest) {
				r.Title = testutils.StringPtr("Updated Title")
				r.Description = testutils.StringPtr("Updated Description")
				r.Status = testutils.TaskStatusPtr(model.TaskStatusCompleted)
			}),
			wantErr: false,
		},
		{
			name:    "non_existing_task",
			id:      999,
			userID:  userID,
			req:     testutils.UpdateTaskRequestFixture(),
			wantErr: true,
		},
		{
			name:    "wrong_user",
			id:      createdTask.ID,
			userID:  999,
			req:     testutils.UpdateTaskRequestFixture(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			task, err := suite.repo.Update(suite.ctx, tt.id, tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), task)
				return
			}

			require.NoError(suite.T(), err)
			require.NotNil(suite.T(), task)

			assert.Equal(suite.T(), tt.id, task.ID)
			assert.Equal(suite.T(), tt.userID, task.UserID)

			if tt.req.Title != nil {
				assert.Equal(suite.T(), *tt.req.Title, task.Title)
			}
			if tt.req.Description != nil {
				assert.Equal(suite.T(), *tt.req.Description, task.Description)
			}
			if tt.req.Status != nil {
				assert.Equal(suite.T(), *tt.req.Status, task.Status)
			}

			assert.True(suite.T(), task.UpdatedAt.After(task.CreatedAt))
		})
	}
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
