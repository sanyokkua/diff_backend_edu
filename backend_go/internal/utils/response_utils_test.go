package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go_backend/internal/apperrors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormatErrorMessage(t *testing.T) {
	// Test basic error formatting
	err := errors.New("test error")
	formattedMessage := formatErrorMessage(err)
	assert.Equal(t, "*errors.errorString: test error", formattedMessage)

	err = apperrors.NewAccessDeniedError("test error")
	formattedMessage = formatErrorMessage(err)
	assert.Equal(t, "apperrors.AccessDeniedError: test error", formattedMessage)
}

func TestCreateErrorResponse(t *testing.T) {
	err := errors.New("sample error")
	statusCode := http.StatusBadRequest

	response := createErrorResponse(statusCode, err)

	assert.Equal(t, statusCode, response.StatusCode)
	assert.Equal(t, http.StatusText(statusCode), response.StatusMessage)
	assert.Contains(t, response.Error, "sample error")
}

func TestCreateResponse(t *testing.T) {
	data := "test data"
	statusCode := http.StatusOK

	response := createResponse(data, statusCode)

	assert.Equal(t, statusCode, response.StatusCode)
	assert.Equal(t, data, response.Data)
	assert.Equal(t, http.StatusText(statusCode), response.StatusMessage)
	assert.Empty(t, response.Error)
}

func TestWriteSuccessResponse(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Test data
	data := "success data"
	statusCode := http.StatusOK

	// Call function
	WriteSuccessResponse(ctx, statusCode, data)

	// Verify response
	assert.Equal(t, statusCode, w.Code)
	assert.Contains(t, w.Body.String(), "success data")
}

func TestWriteErrorResponse(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Test error
	err := errors.New("test error")
	statusCode := http.StatusBadRequest

	// Call function
	WriteErrorResponse(ctx, statusCode, err)

	// Verify response
	assert.Equal(t, statusCode, w.Code)
	assert.Contains(t, w.Body.String(), "test error")
}

func TestGetErrStatusCode(t *testing.T) {
	// Test EmailAlreadyExistsError
	err1 := apperrors.EmailAlreadyExistsError{}
	statusCode := GetErrStatusCode(err1, 0)
	assert.Equal(t, http.StatusConflict, statusCode)

	// Test InvalidEmailFormatError
	err2 := apperrors.InvalidEmailFormatError{}
	statusCode = GetErrStatusCode(err2, 0)
	assert.Equal(t, http.StatusBadRequest, statusCode)

	// Test InvalidJwtTokenError
	err3 := apperrors.InvalidJwtTokenError{}
	statusCode = GetErrStatusCode(err3, 0)
	assert.Equal(t, http.StatusBadRequest, statusCode)

	// Test InvalidPasswordError
	err4 := apperrors.InvalidPasswordError{}
	statusCode = GetErrStatusCode(err4, 0)
	assert.Equal(t, http.StatusBadRequest, statusCode)

	// Test TaskAlreadyExistsError
	err5 := apperrors.TaskAlreadyExistsError{}
	statusCode = GetErrStatusCode(err5, 0)
	assert.Equal(t, http.StatusConflict, statusCode)

	// Test TaskNotFoundError
	err6 := apperrors.TaskNotFoundError{}
	statusCode = GetErrStatusCode(err6, 0)
	assert.Equal(t, http.StatusNotFound, statusCode)

	// Test AccessDeniedError
	err7 := apperrors.AccessDeniedError{}
	statusCode = GetErrStatusCode(err7, 0)
	assert.Equal(t, http.StatusForbidden, statusCode)

	// Test AuthenticationCredentialsNotFoundError
	err8 := apperrors.AuthenticationCredentialsNotFoundError{}
	statusCode = GetErrStatusCode(err8, 0)
	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// Test InsufficientAuthenticationError
	err9 := apperrors.InsufficientAuthenticationError{}
	statusCode = GetErrStatusCode(err9, 0)
	assert.Equal(t, http.StatusUnauthorized, statusCode)

	// Test NoHandlerFoundError
	err10 := apperrors.NoHandlerFoundError{}
	statusCode = GetErrStatusCode(err10, 0)
	assert.Equal(t, http.StatusNotFound, statusCode)

	// Test IllegalArgumentError
	err11 := apperrors.IllegalArgumentError{}
	statusCode = GetErrStatusCode(err11, 0)
	assert.Equal(t, http.StatusBadRequest, statusCode)

	// Test default case with a non-app-specific error
	err12 := errors.New("generic error")
	statusCode = GetErrStatusCode(err12, 0)
	assert.Equal(t, http.StatusInternalServerError, statusCode)

	// Test default case with a non-app-specific error and default status set to BadRequest
	err13 := errors.New("generic error")
	statusCode = GetErrStatusCode(err13, http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
