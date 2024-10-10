package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_backend/dto"
	"go_backend/models"
	"gorm.io/gorm"
	"testing"
)

func TestCreateTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockUserRepo := new(MockUserRepository)

	taskService := NewTaskService(mockTaskRepo, mockUserRepo)

	userId := int64(1)
	taskCreationDTO := &dto.TaskCreationDTO{
		Name:        "Test Task",
		Description: "Task Description",
	}

	// User to be returned from the mockUserRepo
	user := &models.User{
		UserID: userId,
		Email:  "test@example.com",
	}

	// Task to be created
	newTask := &models.Task{
		Name:        taskCreationDTO.Name,
		Description: taskCreationDTO.Description,
		UserID:      user.UserID,
	}

	var taskNil *models.Task = nil
	// Set up the mocks
	mockUserRepo.On("GetUserByID", userId).Return(user, nil)
	mockTaskRepo.On("FindByUserAndName", user, taskCreationDTO.Name).Return(taskNil, gorm.ErrRecordNotFound) // No existing task
	mockTaskRepo.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(newTask, nil)

	// Execute the service method
	result, err := taskService.CreateTask(userId, taskCreationDTO)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, taskCreationDTO.Name, result.Name)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockUserRepo := new(MockUserRepository)

	taskService := NewTaskService(mockTaskRepo, mockUserRepo)

	userId := int64(1)
	taskId := int64(1)
	taskUpdateDTO := &dto.TaskUpdateDTO{
		Name:        "Updated Task",
		Description: "Updated Description",
	}

	// Existing task and user
	existingTask := &models.Task{
		TaskID:      taskId,
		Name:        "Old Task",
		Description: "Old Description",
		UserID:      userId,
	}

	// Set up the mocks
	mockTaskRepo.On("FindByTaskID", taskId).Return(existingTask, nil)
	mockTaskRepo.On("UpdateTask", mock.AnythingOfType("*models.Task")).Return(existingTask, nil)

	// Execute the service method
	updatedTask, err := taskService.UpdateTask(userId, taskId, taskUpdateDTO)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, taskUpdateDTO.Name, updatedTask.Name)
	assert.Equal(t, taskUpdateDTO.Description, updatedTask.Description)

	mockTaskRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockUserRepo := new(MockUserRepository)

	taskService := NewTaskService(mockTaskRepo, mockUserRepo)

	userId := int64(1)
	taskId := int64(1)

	// Task and user to be fetched
	existingTask := &models.Task{
		TaskID: taskId,
		UserID: userId,
	}

	// Set up the mocks
	mockTaskRepo.On("FindByTaskID", taskId).Return(existingTask, nil)
	mockTaskRepo.On("DeleteTask", existingTask).Return(nil)

	// Execute the service method
	err := taskService.DeleteTask(userId, taskId)

	// Assert results
	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAllTasksForUser(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockUserRepo := new(MockUserRepository)

	taskService := NewTaskService(mockTaskRepo, mockUserRepo)

	userId := int64(1)

	// Mock tasks and user
	tasks := []models.Task{
		{Name: "Task 1", Description: "Description 1", UserID: userId},
		{Name: "Task 2", Description: "Description 2", UserID: userId},
	}

	user := &models.User{
		UserID: userId,
		Email:  "test@example.com",
	}

	// Set up the mocks
	mockUserRepo.On("GetUserByID", userId).Return(user, nil)
	mockTaskRepo.On("FindAllByUser", user).Return(tasks, nil)

	// Execute the service method
	result, err := taskService.GetAllTasksForUser(userId)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockTaskRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetTaskByUserIDAndTaskID(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockUserRepo := new(MockUserRepository)

	taskService := NewTaskService(mockTaskRepo, mockUserRepo)

	userId := int64(1)
	taskId := int64(1)

	// Existing task and user
	existingTask := &models.Task{
		TaskID:      taskId,
		Name:        "Task 1",
		Description: "Description 1",
		UserID:      userId,
	}

	user := &models.User{
		UserID: userId,
		Email:  "test@example.com",
	}

	// Set up the mocks
	mockUserRepo.On("GetUserByID", userId).Return(user, nil)
	mockTaskRepo.On("FindByUserAndTaskID", user, taskId).Return(existingTask, nil)

	// Execute the service method
	result, err := taskService.GetTaskByUserIDAndTaskID(userId, taskId)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, existingTask.Name, result.Name)

	mockTaskRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
