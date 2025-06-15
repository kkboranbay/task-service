package testutils

import (
	"github.com/kkboranbay/task-service/internal/model"
	"time"
)

func TaskFixture(overrides ...func(*model.Task)) *model.Task {
	task := &model.Task{
		ID:          1,
		UserID:      1,
		Title:       "Test task",
		Description: "Test Description",
		Status:      model.TaskStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, override := range overrides {
		override(task)
	}

	return task
}

func CreateTaskRequestFixture(overrides ...func(*model.CreateTaskRequest)) model.CreateTaskRequest {
	req := model.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
	}

	for _, override := range overrides {
		override(&req)
	}

	return req
}

func UpdateTaskRequestFixture(overrides ...func(*model.UpdateTaskRequest)) model.UpdateTaskRequest {
	status := model.TaskStatusInProgress
	req := model.UpdateTaskRequest{
		Title:       StringPtr("Updated Test Task"),
		Description: StringPtr("Updated Test Description"),
		Status:      &status,
	}

	for _, override := range overrides {
		override(&req)
	}

	return req
}

func LoginRequestFixture(overrides ...func(*model.LoginRequest)) model.LoginRequest {
	req := model.LoginRequest{
		Username: "admin",
		Password: "admin",
	}

	for _, override := range overrides {
		override(&req)
	}

	return req
}

func StringPtr(s string) *string {
	return &s
}
