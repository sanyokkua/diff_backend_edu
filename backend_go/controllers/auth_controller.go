package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/dto"
	"go_backend/utils"
	"net/http"
)

// AuthController handles authentication-related requests, including user login and registration.
type AuthController struct {
	authService api.AuthenticationService // The service used for user authentication
}

// NewAuthController creates a new instance of AuthController.
//
// Parameters:
//   - authService: an instance of AuthenticationService used for handling user authentication.
//
// Returns:
//   - *AuthController: a pointer to the newly created AuthController instance.
func NewAuthController(authService api.AuthenticationService) *AuthController {
	log.Debug().Msg("Creating new AuthController")
	return &AuthController{authService: authService}
}

// RegisterAuthRoutes registers authentication-related routes with the provided router.
//
// Parameters:
//   - router: the Gin router to which the authentication routes will be registered.
//   - authController: the AuthController instance that handles authentication logic.
func RegisterAuthRoutes(router *gin.Engine, authController *AuthController) {
	log.Debug().Msg("Registering authentication routes")

	v1 := router.Group("api/v1/auth") // Create a new route group for API version 1 authentication
	{
		v1.POST("/login", authController.loginUser)       // Route for user login
		v1.POST("/register", authController.registerUser) // Route for user registration
	}
}

// loginUser handles user login requests.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function attempts to bind the incoming JSON request body to a UserLoginDTO,
// validates the data, and calls the authentication service to log the user in.
// If successful, it returns user details; otherwise, it responds with an error.
func (a *AuthController) loginUser(ctx *gin.Context) {
	log.Debug().Msg("Handling login request")

	var loginDTO dto.UserLoginDTO
	// Bind JSON request to loginDTO
	if err := ctx.BindJSON(&loginDTO); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON for login")
		utils.WriteErrorResponse(ctx, utils.GetErrStatusCode(err, http.StatusBadRequest), err)
		return
	}

	// Attempt to log the user in using the authentication service
	userDTO, err := a.authService.LoginUser(&loginDTO)
	if err != nil {
		log.Error().Err(err).Msg("Failed to login user")
		utils.WriteErrorResponse(ctx, utils.GetErrStatusCode(err, http.StatusBadRequest), err)
		return
	}

	log.Debug().Str("email", userDTO.Email).Msg("User logged in successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, userDTO) // Return success response with user details
}

// registerUser handles user registration requests.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function attempts to bind the incoming JSON request body to a UserCreationDTO,
// validates the data, and calls the authentication service to register the user.
// If successful, it returns the created user details; otherwise, it responds with an error.
func (a *AuthController) registerUser(ctx *gin.Context) {
	log.Debug().Msg("Handling registration request")

	var registerDTO dto.UserCreationDTO
	// Bind JSON request to registerDTO
	if err := ctx.BindJSON(&registerDTO); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON for registration")
		utils.WriteErrorResponse(ctx, utils.GetErrStatusCode(err, http.StatusBadRequest), err)
		return
	}

	// Attempt to register the user using the authentication service
	userDTO, err := a.authService.RegisterUser(&registerDTO)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		utils.WriteErrorResponse(ctx, utils.GetErrStatusCode(err, http.StatusBadRequest), err)
		return
	}

	log.Debug().Str("email", userDTO.Email).Msg("User registered successfully")
	utils.WriteSuccessResponse(ctx, http.StatusCreated, userDTO) // Return success response with user details
}
