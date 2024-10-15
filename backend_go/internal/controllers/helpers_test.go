package controllers

import (
	"github.com/stretchr/testify/mock"
	"go_backend/internal/dto"
)

// MockUserService is a mock implementation of the UserService interface.
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(userCreationDTO *dto.UserCreationDTO) (*dto.UserDTO, error) {
	args := m.Called(userCreationDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdatePassword(userID int64, userUpdateDTO *dto.UserUpdateDTO) (*dto.UserDTO, error) {
	args := m.Called(userID, userUpdateDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func (m *MockUserService) Delete(userID int64, userDeletionDTO *dto.UserDeletionDTO) error {
	args := m.Called(userID, userDeletionDTO)
	return args.Error(0)
}

// MockAuthService is a mock implementation of the AuthenticationService interface.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) LoginUser(loginDTO *dto.UserLoginDTO) (*dto.UserDTO, error) {
	args := m.Called(loginDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func (m *MockAuthService) RegisterUser(registerDTO *dto.UserCreationDTO) (*dto.UserDTO, error) {
	args := m.Called(registerDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

// MockTaskService is a mock implementation of the TaskService interface.
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(userID int64, taskCreationDTO *dto.TaskCreationDTO) (*dto.TaskDTO, error) {
	args := m.Called(userID, taskCreationDTO)
	return args.Get(0).(*dto.TaskDTO), args.Error(1)
}

func (m *MockTaskService) GetTaskByUserIDAndTaskID(userID, taskID int64) (*dto.TaskDTO, error) {
	args := m.Called(userID, taskID)
	return args.Get(0).(*dto.TaskDTO), args.Error(1)
}

func (m *MockTaskService) UpdateTask(userID, taskID int64, taskUpdateDTO *dto.TaskUpdateDTO) (*dto.TaskDTO, error) {
	args := m.Called(userID, taskID, taskUpdateDTO)
	return args.Get(0).(*dto.TaskDTO), args.Error(1)
}

func (m *MockTaskService) DeleteTask(userID, taskID int64) error {
	args := m.Called(userID, taskID)
	return args.Error(0)
}

func (m *MockTaskService) GetAllTasksForUser(userID int64) ([]*dto.TaskDTO, error) {
	args := m.Called(userID)
	return args.Get(0).([]*dto.TaskDTO), args.Error(1)
}
