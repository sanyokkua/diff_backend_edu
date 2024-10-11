package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_backend/apperrors"
	"go_backend/dto"
	"go_backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper to set up Gin, MockService, and routes
func setupTaskTestEnvironment() (*gin.Engine, *MockTaskService) {
	gin.SetMode(gin.TestMode)
	mockTaskService := new(MockTaskService)
	taskController := NewTaskController(mockTaskService)

	r := gin.Default()
	RegisterTaskRoutes(r, taskController, func(ctx *gin.Context) {
		ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
		ctx.Next()
	}) // mock middleware

	return r, mockTaskService
}

// Test for createTask
func TestCreateTask(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		body           *dto.TaskCreationDTO
		mockService    func(service *MockTaskService)
		expectedStatus int
	}{
		{
			name:   "Valid task creation",
			userID: 1,
			body: &dto.TaskCreationDTO{
				Name:        "New Task",
				Description: "Task description",
			},
			mockService: func(service *MockTaskService) {
				service.On("CreateTask", int64(1), mock.Anything).Return(&dto.TaskDTO{TaskID: 1, Name: "New Task", Description: "Task description", UserID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request body",
			userID:         1,
			body:           &dto.TaskCreationDTO{Name: ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockTaskService := setupTaskTestEnvironment()

			if tt.mockService != nil {
				tt.mockService(mockTaskService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/1/tasks/", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockTaskService.AssertExpectations(t)
			}
		})
	}
}

// Test for getTaskByID
func TestGetTaskByID(t *testing.T) {
	var nilTask *dto.TaskDTO = nil
	tests := []struct {
		name           string
		userID         int64
		taskID         int64
		mockService    func(service *MockTaskService)
		expectedStatus int
	}{
		{
			name:   "Valid task retrieval",
			userID: 1,
			taskID: 1,
			mockService: func(service *MockTaskService) {
				service.On("GetTaskByUserIDAndTaskID", int64(1), int64(1)).Return(&dto.TaskDTO{TaskID: 1, Name: "Task 1", Description: "Task description", UserID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Task not found",
			userID: 1,
			taskID: 999,
			mockService: func(service *MockTaskService) {
				service.On("GetTaskByUserIDAndTaskID", int64(1), int64(1)).Return(nilTask, apperrors.TaskNotFoundError{})
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockTaskService := setupTaskTestEnvironment()

			if tt.mockService != nil {
				tt.mockService(mockTaskService)
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/1/tasks/1", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockTaskService.AssertExpectations(t)
			}
		})
	}
}

// Test for updateTask
func TestUpdateTask(t *testing.T) {
	var nilTask *dto.TaskDTO = nil
	tests := []struct {
		name           string
		userID         int64
		taskID         int64
		body           *dto.TaskUpdateDTO
		mockService    func(service *MockTaskService)
		expectedStatus int
	}{
		{
			name:   "Valid task update",
			userID: 1,
			taskID: 1,
			body: &dto.TaskUpdateDTO{
				Name:        "Updated Task",
				Description: "Updated description",
			},
			mockService: func(service *MockTaskService) {
				service.On("UpdateTask", int64(1), int64(1), mock.Anything).Return(&dto.TaskDTO{TaskID: 1, Name: "Updated Task", Description: "Updated description", UserID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid task ID",
			userID: 1,
			taskID: 999,
			body:   &dto.TaskUpdateDTO{Name: "Updated Task"},
			mockService: func(service *MockTaskService) {
				service.On("UpdateTask", int64(1), int64(1), mock.Anything).Return(nilTask, apperrors.TaskNotFoundError{})
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockTaskService := setupTaskTestEnvironment()

			if tt.mockService != nil {
				tt.mockService(mockTaskService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, "/api/v1/users/1/tasks/1", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockTaskService.AssertExpectations(t)
			}
		})
	}
}

// Test for deleteTask
func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		taskID         int64
		mockService    func(service *MockTaskService)
		expectedStatus int
	}{
		{
			name:   "Valid task deletion",
			userID: 1,
			taskID: 1,
			mockService: func(service *MockTaskService) {
				service.On("DeleteTask", int64(1), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Task not found",
			userID: 1,
			taskID: 999,
			mockService: func(service *MockTaskService) {
				service.On("DeleteTask", int64(1), int64(1)).Return(apperrors.TaskNotFoundError{})
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockTaskService := setupTaskTestEnvironment()

			if tt.mockService != nil {
				tt.mockService(mockTaskService)
			}

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/users/1/tasks/1", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockTaskService.AssertExpectations(t)
			}
		})
	}
}

// Test for getAllTasksForUser
func TestGetAllTasksForUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		mockService    func(service *MockTaskService)
		expectedStatus int
	}{
		{
			name:   "Valid task retrieval",
			userID: 1,
			mockService: func(service *MockTaskService) {
				service.On("GetAllTasksForUser", int64(1)).Return([]*dto.TaskDTO{
					{TaskID: 1, Name: "Task 1", Description: "Task 1 description", UserID: 1},
					{TaskID: 2, Name: "Task 2", Description: "Task 2 description", UserID: 1},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "No tasks found",
			userID: 1,
			mockService: func(service *MockTaskService) {
				service.On("GetAllTasksForUser", int64(1)).Return([]*dto.TaskDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockTaskService := setupTaskTestEnvironment()

			if tt.mockService != nil {
				tt.mockService(mockTaskService)
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/1/tasks/", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockTaskService.AssertExpectations(t)
			}
		})
	}
}
