package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
	"gorm.io/gorm"
	"testing"
)

// Mock for TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindByTaskID(taskID int64) (*models.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUserAndName(user *models.User, name string) (*models.Task, error) {
	args := m.Called(user, name)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUserAndTaskID(user *models.User, taskID int64) (*models.Task, error) {
	args := m.Called(user, taskID)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAllByUser(user *models.User) ([]models.Task, error) {
	args := m.Called(user)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func TestCheckTaskExistsForUser_TaskNotFound(t *testing.T) {
	// Mock setup
	mockRepo := new(MockTaskRepository)
	user := &models.User{UserID: 1}
	taskName := "testTask"
	var task models.Task

	mockRepo.On("FindByUserAndName", user, taskName).Return(&task, gorm.ErrRecordNotFound)

	// Call the function
	err := CheckTaskExistsForUser(mockRepo, user, taskName)

	// Validate the result
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCheckTaskExistsForUser_TaskExists(t *testing.T) {
	// Mock setup
	mockRepo := new(MockTaskRepository)
	user := &models.User{UserID: 1}
	taskName := "testTask"

	mockRepo.On("FindByUserAndName", user, taskName).Return(&models.Task{}, nil)

	// Call the function
	err := CheckTaskExistsForUser(mockRepo, user, taskName)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.TaskAlreadyExistsError{}, err)
	mockRepo.AssertExpectations(t)
}

func TestCheckTaskExistsForUser_RepoError(t *testing.T) {
	// Mock setup
	mockRepo := new(MockTaskRepository)
	user := &models.User{UserID: 1}
	taskName := "testTask"
	var task models.Task

	mockRepo.On("FindByUserAndName", user, taskName).Return(&task, errors.New("db error"))

	// Call the function
	err := CheckTaskExistsForUser(mockRepo, user, taskName)

	// Validate the result
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestValidateTaskCreation_NilDTO(t *testing.T) {
	// Call the function
	err := ValidateTaskCreation(nil)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "TaskCreationDTO is nil", err.Error())
}

func TestValidateTaskCreation_EmptyName(t *testing.T) {
	taskDTO := &dto.TaskCreationDTO{Name: "", Description: "Description"}

	// Call the function
	err := ValidateTaskCreation(taskDTO)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "Task name cannot be empty", err.Error())
}

func TestValidateTaskCreation_EmptyDescription(t *testing.T) {
	taskDTO := &dto.TaskCreationDTO{Name: "TaskName", Description: ""}

	// Call the function
	err := ValidateTaskCreation(taskDTO)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "Task description cannot be empty", err.Error())
}

func TestValidateTaskCreation_ValidDTO(t *testing.T) {
	taskDTO := &dto.TaskCreationDTO{Name: "TaskName", Description: "Description"}

	// Call the function
	err := ValidateTaskCreation(taskDTO)

	// Validate the result
	assert.NoError(t, err)
}

func TestValidateTaskUpdateDTO_NilDTO(t *testing.T) {
	// Call the function
	err := ValidateTaskUpdateDTO(nil)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "TaskUpdateDTO is nil", err.Error())
}

func TestValidateTaskUpdateDTO_EmptyName(t *testing.T) {
	taskDTO := &dto.TaskUpdateDTO{Name: "", Description: "Description"}

	// Call the function
	err := ValidateTaskUpdateDTO(taskDTO)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "Task name cannot be empty", err.Error())
}

func TestValidateTaskUpdateDTO_EmptyDescription(t *testing.T) {
	taskDTO := &dto.TaskUpdateDTO{Name: "TaskName", Description: ""}

	// Call the function
	err := ValidateTaskUpdateDTO(taskDTO)

	// Validate the result
	assert.Error(t, err)
	assert.IsType(t, apperrors.IllegalArgumentError{}, err)
	assert.Equal(t, "Task description cannot be empty", err.Error())
}

func TestValidateTaskUpdateDTO_ValidDTO(t *testing.T) {
	taskDTO := &dto.TaskUpdateDTO{Name: "TaskName", Description: "Description"}

	// Call the function
	err := ValidateTaskUpdateDTO(taskDTO)

	// Validate the result
	assert.NoError(t, err)
}
