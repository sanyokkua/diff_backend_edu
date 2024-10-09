package services

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/apperrors"
	"go_backend/dto"
	"go_backend/utils"
)

// authenticationService is a struct that implements the AuthenticationService interface.
// It provides methods for user authentication and registration.
type authenticationService struct {
	userService     api.UserService     // Service for user-related operations.
	userRepository  api.UserRepository  // Repository for accessing user data.
	jwtService      api.JwtService      // Service for generating and validating JWT tokens.
	passwordEncoder api.PasswordEncoder // Encoder for hashing and verifying passwords.
}

// NewAuthenticationService creates a new AuthenticationService with the specified dependencies.
//
// Parameters:
//   - userService: An instance of UserService for handling user operations.
//   - userRepository: An instance of UserRepository for accessing user data.
//   - jwtService: An instance of JwtService for JWT operations.
//   - passwordEncoder: An instance of PasswordEncoder for password operations.
//
// Returns:
//   - api.AuthenticationService: A new instance of AuthenticationService.
func NewAuthenticationService(userService api.UserService, userRepository api.UserRepository, jwtService api.JwtService, passwordEncoder api.PasswordEncoder) api.AuthenticationService {
	log.Debug().Msg("Creating new AuthenticationService")
	return &authenticationService{
		userService:     userService,
		userRepository:  userRepository,
		jwtService:      jwtService,
		passwordEncoder: passwordEncoder,
	}
}

// LoginUser validates the login credentials and generates a JWT token for the user.
//
// Parameters:
//   - loginDto: A pointer to UserLoginDTO containing login credentials.
//
// Returns:
//   - *dto.UserDTO: A pointer to UserDTO containing user information and JWT token.
//   - error: An error if the login fails.
func (a *authenticationService) LoginUser(loginDto *dto.UserLoginDTO) (*dto.UserDTO, error) {
	log.Debug().Str("email", loginDto.Email).Msg("Validating login credentials")

	// Validate login DTO
	if err := utils.ValidateUserLoginDto(loginDto); err != nil {
		log.Error().Err(err).Msg("Validation failed for UserLoginDTO")
		return nil, err
	}

	// Fetch user by email
	user, err := a.userRepository.GetUserByEmail(loginDto.Email)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching user by email")
		return nil, apperrors.NewIllegalArgumentError(fmt.Sprintf("failed to retrieve user: %v", err))
	}

	// Validate password
	matches, err := a.passwordEncoder.Matches(loginDto.Password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("Password matching failed")
		return nil, apperrors.NewInvalidPasswordError(err.Error())
	}

	if !matches {
		log.Warn().Msg("Invalid credentials")
		return nil, apperrors.NewInvalidPasswordError("Invalid credentials")
	}

	// Generate JWT token
	jwtToken, err := a.jwtService.GenerateJwtToken(user.Email)
	if err != nil {
		log.Error().Err(err).Msg("JWT token generation failed")
		return nil, apperrors.NewInvalidJwtTokenError(err.Error())
	}

	// Return UserDTO with the generated token
	userDto := &dto.UserDTO{
		UserID:   user.UserID,
		Email:    user.Email,
		JwtToken: jwtToken,
	}
	log.Debug().Int64("userID", user.UserID).Msg("User logged in successfully")
	return userDto, nil
}

// RegisterUser validates and registers a new user, then generates a JWT token.
//
// Parameters:
//   - userCreationDTO: A pointer to UserCreationDTO containing registration details.
//
// Returns:
//   - *dto.UserDTO: A pointer to UserDTO containing the newly registered user information and JWT token.
//   - error: An error if the registration fails.
func (a *authenticationService) RegisterUser(userCreationDTO *dto.UserCreationDTO) (*dto.UserDTO, error) {
	log.Debug().Str("email", userCreationDTO.Email).Msg("Validating registration details")

	// Validate user creation DTO
	if err := utils.ValidateUserCreationDTO(userCreationDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed for UserCreationDTO")
		return nil, err
	}

	// Create the new user
	newUser, err := a.userService.Create(userCreationDTO)
	if err != nil {
		log.Error().Err(err).Msg("User creation failed")
		return nil, apperrors.NewIllegalArgumentError(err.Error())
	}

	// Generate JWT token for the new user
	jwtToken, err := a.jwtService.GenerateJwtToken(newUser.Email)
	if err != nil {
		log.Error().Err(err).Msg("JWT token generation failed")
		return nil, apperrors.NewInvalidJwtTokenError(err.Error())
	}

	// Return the UserDTO with the generated token
	userDto := &dto.UserDTO{
		UserID:   newUser.UserID,
		Email:    newUser.Email,
		JwtToken: jwtToken,
	}
	log.Debug().Int64("userID", newUser.UserID).Msg("User registered successfully")
	return userDto, nil
}
