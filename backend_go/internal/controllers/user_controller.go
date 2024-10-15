package controllers

import (
	"github.com/rs/zerolog/log"
	"go_backend/internal/api"
	"go_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_backend/internal/dto"
	"go_backend/internal/utils"
)

// UserController handles user-related requests.
type UserController struct {
	userService api.UserService // Service to handle user-related operations
}

// NewUserController creates a new instance of UserController.
//
// Parameters:
//   - userService: an instance of UserService for managing user operations.
//
// Returns:
//   - *UserController: a pointer to the newly created UserController instance.
func NewUserController(userService api.UserService) *UserController {
	log.Debug().Msg("Initializing new UserController")
	return &UserController{
		userService: userService,
	}
}

// RegisterUserRoutes registers user routes with the provided router.
//
// Parameters:
//   - router: the Gin router to which the user routes will be registered.
//   - userController: the UserController instance that handles user logic.
//   - middleware: a middleware function to apply to the user routes.
func RegisterUserRoutes(router *gin.Engine, userController *UserController, middleware gin.HandlerFunc) {
	v1 := router.Group("api/v1/users") // Grouping routes under user operations
	v1.Use(middleware)                 // Apply middleware for authentication
	{
		v1.GET("/:userId", userController.getUserByID)                 // Route to get user by ID
		v1.PUT("/:userId/password", userController.updateUserPassword) // Route to update user password
		v1.POST("/:userId/delete", userController.deleteUser)          // Route to delete a user
	}
}

// extractAndValidateUser handles extracting and validating the user ID from the request.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// Returns:
//   - userID: the ID of the user extracted from the request context.
//   - userFromContext: the user object extracted from the context.
//   - error: any error encountered during extraction and validation.
func (r *UserController) extractAndValidateUser(ctx *gin.Context) (int64, *models.User, error) {
	userFromContext, err := getUserFromContext(ctx) // Get the authenticated user from context
	if err != nil {
		return 0, nil, err
	}

	userID, err := getIntParamByName(ctx, "userId") // Extract user ID from request parameters
	if err != nil {
		return 0, nil, err
	}

	if err := utils.ValidateAuthenticatedUserID(userFromContext.UserID, userID); err != nil {
		return 0, nil, err // Validate that the authenticated user can access the requested user ID
	}

	return userID, userFromContext, nil
}

// getUserByID handles requests to retrieve a user by ID.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user ID, retrieves the user from the service,
// and responds with the user information or an error if encountered.
func (r *UserController) getUserByID(ctx *gin.Context) {
	log.Debug().Msg("Handling getUserByID request")

	userID, userFromContext, err := r.extractAndValidateUser(ctx) // Extract and validate user ID
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	userDTO := dto.UserDTO{
		UserID: userFromContext.UserID, // Create UserDTO from the context user
		Email:  userFromContext.Email,
	}

	log.Info().Int64("userId", userID).Msg("User retrieved successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, userDTO) // Return user details
}

// updateUserPassword handles requests to update a user's password.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user ID, binds the request body to a UserUpdateDTO,
// and calls the user service to update the user's password. It responds with
// the updated password or an error if encountered.
func (r *UserController) updateUserPassword(ctx *gin.Context) {
	log.Debug().Msg("Handling updateUserPassword request")

	userID, _, err := r.extractAndValidateUser(ctx) // Extract and validate user ID
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var userUpdateDTO dto.UserUpdateDTO
	if err := ctx.BindJSON(&userUpdateDTO); err != nil { // Bind incoming JSON to DTO
		log.Error().Err(err).Msg("Failed to parse UserUpdateDTO")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	updatedPassword, err := r.userService.UpdatePassword(userID, &userUpdateDTO) // Update password via user service
	if err != nil {
		log.Error().Err(err).Int64("userId", userID).Msg("Failed to update password")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("userId", userID).Msg("Password updated successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, updatedPassword) // Return updated password information
}

// deleteUser handles requests to delete a user.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user ID, binds the request body to a UserDeletionDTO,
// and calls the user service to delete the user. It responds with a no content status
// or an error if encountered.
func (r *UserController) deleteUser(ctx *gin.Context) {
	log.Debug().Msg("Handling deleteUser request")

	userID, _, err := r.extractAndValidateUser(ctx) // Extract and validate user ID
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var userDeletionDTO dto.UserDeletionDTO
	if err := ctx.BindJSON(&userDeletionDTO); err != nil { // Bind incoming JSON to DTO
		log.Error().Err(err).Msg("Failed to parse UserDeletionDTO")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := r.userService.Delete(userID, &userDeletionDTO); err != nil { // Call user service to delete user
		log.Error().Err(err).Int64("userId", userID).Msg("Failed to delete user")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("userId", userID).Msg("User deleted successfully")
	utils.WriteSuccessResponse(ctx, http.StatusNoContent, nil) // Respond with no content on successful deletion
}
