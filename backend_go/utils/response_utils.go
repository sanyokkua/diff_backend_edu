package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/apperrors"
	"go_backend/dto"
	"net/http"
)

// formatErrorMessage formats the error message with its type.
//
// Parameters:
//   - err: An error to format.
//
// Returns:
//   - string: A formatted error message indicating the type and message of the error.
func formatErrorMessage(err error) string {
	log.Debug().Msg("Entering formatErrorMessage")
	defer log.Debug().Msg("Exiting formatErrorMessage")
	return fmt.Sprintf("%T: %s", err, err.Error())
}

// createErrorResponse creates a response body for errors.
//
// Parameters:
//   - statusCode: The HTTP status code for the response.
//   - err: The error for which the response is being created.
//
// Returns:
//   - *dto.ResponseDto[any]: A ResponseDto containing the error details.
func createErrorResponse(statusCode int, err error) *dto.ResponseDto[any] {
	log.Debug().Int("statusCode", statusCode).Str("error", err.Error()).Msg("Entering createErrorResponse")
	defer log.Debug().Int("statusCode", statusCode).Str("error", err.Error()).Msg("Exiting createErrorResponse")
	return dto.NewFullResponseDto(statusCode, http.StatusText(statusCode), "", formatErrorMessage(err))
}

// createResponse creates a standard response DTO.
//
// Parameters:
//   - data: The data to include in the response.
//   - statusCode: The HTTP status code for the response.
//
// Returns:
//   - *dto.ResponseDto[any]: A ResponseDto containing the data.
func createResponse(data interface{}, statusCode int) *dto.ResponseDto[any] {
	log.Debug().Int("statusCode", statusCode).Msg("Entering createResponse")
	defer log.Debug().Int("statusCode", statusCode).Msg("Exiting createResponse")
	return dto.NewFullResponseDto(statusCode, http.StatusText(statusCode), data, "")
}

// WriteSuccessResponse writes a success response to the context.
//
// Parameters:
//   - ctx: The Gin context to write the response to.
//   - statusCode: The HTTP status code for the response.
//   - data: The data to include in the success response.
func WriteSuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	log.Debug().Int("statusCode", statusCode).Msg("Entering WriteSuccessResponse")
	defer log.Debug().Int("statusCode", statusCode).Msg("Exiting WriteSuccessResponse")
	ctx.JSON(statusCode, createResponse(data, statusCode))
}

// WriteErrorResponse writes an error response to the context.
//
// Parameters:
//   - ctx: The Gin context to write the response to.
//   - statusCode: The HTTP status code for the response.
//   - err: The error to include in the response.
func WriteErrorResponse(ctx *gin.Context, statusCode int, err error) {
	log.Debug().Int("statusCode", statusCode).Str("error", err.Error()).Msg("Entering WriteErrorResponse")
	defer log.Debug().Int("statusCode", statusCode).Str("error", err.Error()).Msg("Exiting WriteErrorResponse")
	ctx.JSON(statusCode, createErrorResponse(statusCode, err))
}

// GetErrStatusCode determines the HTTP status code based on the error type.
//
// Parameters:
//   - err: The error for which to determine the status code.
//   - statusDefault: A default HTTP status code to return if no specific code is found.
//
// Returns:
//   - int: The HTTP status code corresponding to the error type or the default code if no specific match is found.
func GetErrStatusCode(err error, statusDefault int) int {
	log.Debug().Msg("Entering GetErrStatusCode")
	defer log.Debug().Msg("Exiting GetErrStatusCode")

	switch {
	case errors.As(err, &apperrors.EmailAlreadyExistsError{}):
		log.Debug().Int("StatusCode", http.StatusConflict).Str("Error", "EmailAlreadyExistsError").Msg("Status chosen")
		return http.StatusConflict
	case errors.As(err, &apperrors.InvalidEmailFormatError{}):
		log.Debug().Int("StatusCode", http.StatusBadRequest).Str("Error", "InvalidEmailFormatError").Msg("Status chosen")
		return http.StatusBadRequest
	case errors.As(err, &apperrors.InvalidJwtTokenError{}):
		log.Debug().Int("StatusCode", http.StatusBadRequest).Str("Error", "InvalidJwtTokenError").Msg("Status chosen")
		return http.StatusBadRequest
	case errors.As(err, &apperrors.InvalidPasswordError{}):
		log.Debug().Int("StatusCode", http.StatusBadRequest).Str("Error", "InvalidPasswordError").Msg("Status chosen")
		return http.StatusBadRequest
	case errors.As(err, &apperrors.TaskAlreadyExistsError{}):
		log.Debug().Int("StatusCode", http.StatusConflict).Str("Error", "TaskAlreadyExistsError").Msg("Status chosen")
		return http.StatusConflict
	case errors.As(err, &apperrors.TaskNotFoundError{}):
		log.Debug().Int("StatusCode", http.StatusNotFound).Str("Error", "TaskNotFoundError").Msg("Status chosen")
		return http.StatusNotFound
	case errors.As(err, &apperrors.AccessDeniedError{}):
		log.Debug().Int("StatusCode", http.StatusForbidden).Str("Error", "AccessDeniedError").Msg("Status chosen")
		return http.StatusForbidden
	case errors.As(err, &apperrors.AuthenticationCredentialsNotFoundError{}):
		log.Debug().Int("StatusCode", http.StatusUnauthorized).Str("Error", "AuthenticationCredentialsNotFoundError").Msg("Status chosen")
		return http.StatusUnauthorized
	case errors.As(err, &apperrors.InsufficientAuthenticationError{}):
		log.Debug().Int("StatusCode", http.StatusUnauthorized).Str("Error", "InsufficientAuthenticationError").Msg("Status chosen")
		return http.StatusUnauthorized
	case errors.As(err, &apperrors.NoHandlerFoundError{}):
		log.Debug().Int("StatusCode", http.StatusNotFound).Str("Error", "NoHandlerFoundError").Msg("Status chosen")
		return http.StatusNotFound
	case errors.As(err, &apperrors.IllegalArgumentError{}):
		log.Debug().Int("StatusCode", http.StatusBadRequest).Str("Error", "IllegalArgumentError").Msg("Status chosen")
		return http.StatusBadRequest
	default:
		if statusDefault == 0 {
			return http.StatusInternalServerError
		}
		return statusDefault
	}
}
