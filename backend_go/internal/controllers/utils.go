package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/internal/apperrors"
	"go_backend/internal/models"
	"go_backend/internal/utils"
	"strconv"
)

// getIntParamByName retrieves an integer parameter from the context by its name.
//
// Parameters:
//   - ctx: the Gin context from which to retrieve the parameter.
//   - paramName: the name of the parameter to retrieve.
//
// Returns:
//   - int64: the parsed integer value of the parameter.
//   - error: an error if the parameter cannot be found or parsed to an integer.
func getIntParamByName(ctx *gin.Context, paramName string) (int64, error) {
	log.Debug().Str("paramName", paramName).Msg("Retrieving integer parameter")
	defer log.Debug().Str("paramName", paramName).Msg("Finished retrieving integer parameter")

	paramStr := ctx.Param(paramName)                    // Retrieve parameter as a string
	paramInt, err := strconv.ParseInt(paramStr, 10, 64) // Convert string to int64
	if err != nil {
		log.Error().Err(err).Str("paramName", paramName).Msg("Failed to parse parameter to int")
		return 0, fmt.Errorf("parameter %s cannot be parsed to int: %w", paramName, err)
	}

	log.Debug().Int64("paramInt", paramInt).Msg("Parameter parsed successfully")
	return paramInt, nil
}

// getUserFromContext retrieves the current user from the Gin context.
//
// Parameters:
//   - ctx: the Gin context from which to retrieve the user.
//
// Returns:
//   - *models.User: a pointer to the current user if found.
//   - error: an error if the user is not found or cannot be converted to a User type.
func getUserFromContext(ctx *gin.Context) (*models.User, error) {
	log.Debug().Msg("Retrieving user from context")
	defer log.Debug().Msg("Finished retrieving user from context")

	currentUser, found := ctx.Get("CurrentUser") // Retrieve current user from context
	if !found {
		log.Warn().Msg("User not logged in")
		return nil, apperrors.NewAuthenticationCredentialsNotFoundError("user not logged in")
	}

	user, ok := currentUser.(*models.User) // Assert the type of the user
	if !ok {
		log.Error().Msg("Failed to convert saved object to user type")
		return nil, apperrors.NewGenericError("cannot convert saved object to user type")
	}

	log.Debug().Str("email", user.Email).Msg("User retrieved successfully")
	return user, nil
}

// handleErrorResponse is a helper function to handle error responses for the API.
//
// Parameters:
//   - ctx: the Gin context for sending the response.
//   - err: the error that occurred, which will determine the status code.
//   - defaultStatusCode: the default status code to use if an appropriate one cannot be determined.
//
// This function writes an error response to the context based on the error and status code.
func handleErrorResponse(ctx *gin.Context, err error, defaultStatusCode int) {
	statusCode := utils.GetErrStatusCode(err, defaultStatusCode) // Determine the appropriate status code
	utils.WriteErrorResponse(ctx, statusCode, err)               // Write the error response to the context
}
