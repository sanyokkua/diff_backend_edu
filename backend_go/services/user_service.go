package services

import (
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/apperrors"
	"go_backend/dto"
	"go_backend/models"
	"go_backend/utils"
)

// userService implements UserService and provides user management functionalities.
type userService struct {
	userRepo        api.UserRepository  // Repository for user data
	passwordEncoder api.PasswordEncoder // Encoder for user passwords
}

// NewUserService creates a new UserService with the given repository and encoder.
//
// Parameters:
//   - repo: An api.UserRepository instance for user operations.
//   - encoder: An api.PasswordEncoder instance for encoding passwords.
//
// Returns:
//   - api.UserService: A new instance of UserService.
func NewUserService(repo api.UserRepository, encoder api.PasswordEncoder) api.UserService {
	log.Debug().Msg("Creating new UserService instance")
	return &userService{
		userRepo:        repo,
		passwordEncoder: encoder,
	}
}

// Create registers a new user by validating input, checking email, and encoding the password.
//
// Parameters:
//   - userCreationDTO: A pointer to dto.UserCreationDTO containing user registration details.
//
// Returns:
//   - *dto.UserDTO: The created user as a UserDTO.
//   - error: An error if user registration fails.
func (s *userService) Create(userCreationDTO *dto.UserCreationDTO) (*dto.UserDTO, error) {
	log.Debug().Msg("Processing user registration")
	defer log.Debug().Msg("Completed user registration processing")

	if err := s.validateUserCreation(userCreationDTO); err != nil {
		return nil, err
	}

	newUser := models.NewFullUser(0, userCreationDTO.Email, s.encodePassword(userCreationDTO.Password), nil)

	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	userDTO := dto.NewFullUserDto(createdUser.UserID, createdUser.Email, "")
	log.Debug().Int64("userID", createdUser.UserID).Msg("User created successfully")
	return userDTO, nil
}

// UpdatePassword updates the user's password.
//
// Parameters:
//   - userID: The ID of the user whose password is to be updated.
//   - userUpdateDTO: A pointer to dto.UserUpdateDTO containing updated password details.
//
// Returns:
//   - *dto.UserDTO: The updated user as a UserDTO.
//   - error: An error if password update fails.
func (s *userService) UpdatePassword(userID int64, userUpdateDTO *dto.UserUpdateDTO) (*dto.UserDTO, error) {
	log.Debug().Int64("userId", userID).Msg("Processing password update")
	defer log.Debug().Int64("userId", userID).Msg("Completed password update processing")

	if err := utils.ValidateUserUpdateDTO(userUpdateDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed for UserUpdateDTO")
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Error().Err(err).Int64("userId", userID).Msg("Failed to fetch user by ID")
		return nil, err
	}

	if err := s.validatePasswordUpdate(userUpdateDTO, user); err != nil {
		return nil, err
	}

	user.PasswordHash = s.encodePassword(userUpdateDTO.NewPassword)

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return nil, err
	}

	userDTO := dto.NewFullUserDto(updatedUser.UserID, updatedUser.Email, "")
	log.Debug().Int64("userID", updatedUser.UserID).Msg("User password updated successfully")
	return userDTO, nil
}

// Delete removes a user from the system.
//
// Parameters:
//   - userID: The ID of the user to be deleted.
//   - userDeletionDTO: A pointer to dto.UserDeletionDTO containing deletion details.
//
// Returns:
//   - error: An error if user deletion fails.
func (s *userService) Delete(userID int64, userDeletionDTO *dto.UserDeletionDTO) error {
	log.Debug().Int64("userId", userID).Msg("Processing user deletion")
	defer log.Debug().Int64("userId", userID).Msg("Completed user deletion processing")

	if err := utils.ValidateUserDeletionDTO(userDeletionDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed for UserDeletionDTO")
		return err
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Error().Err(err).Int64("userId", userID).Msg("Failed to fetch user by ID")
		return err
	}

	if err := s.validateCurrentPassword(userDeletionDTO.CurrentPassword, user.PasswordHash); err != nil {
		return err
	}

	if err := s.userRepo.DeleteUser(userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return apperrors.NewGenericError(err.Error())
	}

	log.Debug().Int64("userId", userID).Msg("User deleted successfully")
	return nil
}

// validateUserCreation validates the input for user registration.
//
// Parameters:
//   - userCreationDTO: A pointer to dto.UserCreationDTO containing user registration details.
//
// Returns:
//   - error: An error if validation fails.
func (s *userService) validateUserCreation(userCreationDTO *dto.UserCreationDTO) error {
	if err := utils.ValidateUserCreationDTO(userCreationDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed for UserCreationDTO")
		return err
	}

	if err := utils.ValidateEmailFormat(userCreationDTO.Email); err != nil {
		log.Error().Err(err).Str("email", userCreationDTO.Email).Msg("Invalid email format")
		return err
	}

	if err := utils.CheckUserExists(s.userRepo, userCreationDTO.Email); err != nil {
		log.Error().Err(err).Str("email", userCreationDTO.Email).Msg("Email already exists")
		return err
	}

	if err := utils.ValidatePasswords(userCreationDTO.Password, userCreationDTO.PasswordConfirmation); err != nil {
		log.Error().Err(err).Msg("Passwords do not match")
		return err
	}

	return nil
}

// encodePassword encodes the user's password.
//
// Parameters:
//   - password: The plaintext password to encode.
//
// Returns:
//   - string: The encoded password.
func (s *userService) encodePassword(password string) string {
	encodedPassword, err := s.passwordEncoder.Encode(password)
	if err != nil {
		log.Error().Err(err).Msg("Password encoding failed")
		return ""
	}
	return encodedPassword
}

// validatePasswordUpdate validates the password update request.
//
// Parameters:
//   - userUpdateDTO: A pointer to dto.UserUpdateDTO containing updated password details.
//   - user: The user whose password is being updated.
//
// Returns:
//   - error: An error if validation fails.
func (s *userService) validatePasswordUpdate(userUpdateDTO *dto.UserUpdateDTO, user *models.User) error {
	if err := utils.ValidatePasswordUpdate(userUpdateDTO, user, s.passwordEncoder); err != nil {
		log.Error().Err(err).Msg("Password validation failed")
		return err
	}
	return nil
}

// validateCurrentPassword validates the current password provided during user deletion.
//
// Parameters:
//   - currentPassword: The current password entered by the user.
//   - storedPassword: The stored password hash of the user.
//
// Returns:
//   - error: An error if the current password does not match.
func (s *userService) validateCurrentPassword(currentPassword, storedPassword string) error {
	matches, err := s.passwordEncoder.Matches(currentPassword, storedPassword)
	if err != nil {
		log.Error().Err(err).Msg("Password matching failed")
		return err
	}

	if !matches {
		log.Warn().Msg("Invalid current password")
		return apperrors.NewInvalidPasswordError("Invalid current password")
	}

	return nil
}
