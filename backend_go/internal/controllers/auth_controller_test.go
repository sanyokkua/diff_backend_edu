package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper to set up Gin, MockService, and routes
func setupAuthTestEnvironment() (*gin.Engine, *MockAuthService) {
	gin.SetMode(gin.TestMode)
	mockAuthService := new(MockAuthService)
	authController := NewAuthController(mockAuthService)

	r := gin.Default()
	RegisterAuthRoutes(r, authController)
	return r, mockAuthService
}

// Table-driven test for loginUser
func TestLoginUser(t *testing.T) {
	var nilUser *dto.UserDTO = nil
	tests := []struct {
		name           string
		body           *dto.UserLoginDTO
		mockService    func(service *MockAuthService)
		expectedStatus int
	}{
		{
			name: "Valid login",
			body: &dto.UserLoginDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockService: func(service *MockAuthService) {
				service.On("LoginUser", mock.Anything).Return(&dto.UserDTO{UserID: 1, Email: "test@example.com", JwtToken: "jwt-token"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid credentials",
			body: &dto.UserLoginDTO{
				Email:    "wrong@example.com",
				Password: "wrongpassword",
			},
			mockService: func(service *MockAuthService) {
				service.On("LoginUser", mock.Anything).Return(nilUser, apperrors.InvalidPasswordError{})
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request body",
			body:           &dto.UserLoginDTO{Email: "invalid-email"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockAuthService := setupAuthTestEnvironment()
			if tt.mockService != nil {
				tt.mockService(mockAuthService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockAuthService.AssertExpectations(t)
			}
		})
	}
}

// Table-driven test for registerUser
func TestRegisterUser(t *testing.T) {
	var nilUser *dto.UserDTO = nil
	tests := []struct {
		name           string
		body           *dto.UserCreationDTO
		mockService    func(service *MockAuthService)
		expectedStatus int
	}{
		{
			name: "Valid registration",
			body: &dto.UserCreationDTO{
				Email:                "newuser@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			mockService: func(service *MockAuthService) {
				service.On("RegisterUser", mock.Anything).Return(&dto.UserDTO{UserID: 1, Email: "newuser@example.com", JwtToken: "jwt-token"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Password confirmation mismatch",
			body: &dto.UserCreationDTO{
				Email:                "newuser@example.com",
				Password:             "password123",
				PasswordConfirmation: "differentpassword",
			},
			mockService: func(service *MockAuthService) {
				service.On("RegisterUser", mock.Anything).Return(nilUser, apperrors.InvalidPasswordError{})
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request body",
			body:           &dto.UserCreationDTO{Email: "invalid-email"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, mockAuthService := setupAuthTestEnvironment()
			if tt.mockService != nil {
				tt.mockService(mockAuthService)
			}

			jsonValue, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.mockService != nil {
				mockAuthService.AssertExpectations(t)
			}
		})
	}
}
