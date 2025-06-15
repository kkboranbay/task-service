package postgres

import (
	"context"
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

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
