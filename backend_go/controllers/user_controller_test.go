package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_backend/dto"
	"go_backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper to set up Gin, MockService, and routes
func setupTestEnvironment(setupContext func(ctx *gin.Context)) (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(MockUserService)
	userController := NewUserController(mockUserService)

	r := gin.Default()
	RegisterUserRoutes(r, userController, setupContext)
	return r, mockUserService
}

// Table-driven test for GetUserByID
func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(ctx *gin.Context)
		expectedStatus int
	}{
		{
			name: "Valid User ID and authenticated",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
				ctx.Next()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No user in context (unauthenticated)",
			setupContext:   func(ctx *gin.Context) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid user ID",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 2, Email: "test@example.com"})
				ctx.Next()
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := setupTestEnvironment(tt.setupContext)
			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Table-driven test for UpdateUserPassword
func TestUpdateUserPassword(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(ctx *gin.Context)
		body           *dto.UserUpdateDTO
		mockService    func(service *MockUserService)
		expectedStatus int
	}{
		{
			name: "Valid password update",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
				ctx.Next()
			},
			body: &dto.UserUpdateDTO{
				CurrentPassword:         "currentPassword",
				NewPassword:             "newPassword123",
				NewPasswordConfirmation: "newPassword123",
			},
			mockService: func(service *MockUserService) {
				service.On("UpdatePassword", int64(1), mock.Anything).Return(&dto.UserDTO{UserID: 1, Email: "test@example.com"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No user in context (unauthenticated)",
			setupContext:   func(ctx *gin.Context) {},
			body:           &dto.UserUpdateDTO{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid request body",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
				ctx.Next()
			},
			body:           &dto.UserUpdateDTO{CurrentPassword: "currentPassword"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockUserService := setupTestEnvironment(tt.setupContext)
			if tt.mockService != nil {
				tt.mockService(mockUserService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, "/api/v1/users/1/password", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}

// Table-driven test for DeleteUser
func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(ctx *gin.Context)
		body           *dto.UserDeletionDTO
		mockService    func(service *MockUserService)
		expectedStatus int
	}{
		{
			name: "Valid deletion",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
				ctx.Next()
			},
			body: &dto.UserDeletionDTO{
				Email:           "test@example.com",
				CurrentPassword: "currentPassword",
			},
			mockService: func(service *MockUserService) {
				service.On("Delete", int64(1), mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "No user in context (unauthenticated)",
			setupContext:   func(ctx *gin.Context) {},
			body:           &dto.UserDeletionDTO{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid request body",
			setupContext: func(ctx *gin.Context) {
				ctx.Set("CurrentUser", &models.User{UserID: 1, Email: "test@example.com"})
				ctx.Next()
			},
			body:           &dto.UserDeletionDTO{Email: "test@example.com"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockUserService := setupTestEnvironment(tt.setupContext)
			if tt.mockService != nil {
				tt.mockService(mockUserService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/1/delete", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockUserService.AssertExpectations(t)
			}
		})
	}
}
