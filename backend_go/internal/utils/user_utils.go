package utils

import (
	"errors"
	"github.com/rs/zerolog/log"
	"go_backend/internal/api"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
	"gorm.io/gorm"
	"regexp"
)

const EmailPattern = `^[\w-.]+@([\w-]+\.)[\w-]{2,4}$`

var pattern = regexp.MustCompile(EmailPattern)

// ValidateEmailFormat checks if the provided email string is in a valid format.
//
// Parameters:
//   - email: The email string to validate.
//
// Returns:
//   - error: Returns nil if the email format is valid; otherwise, returns an InvalidEmailFormatError.
func ValidateEmailFormat(email string) error {
	log.Debug().Str("email", email).Msg("Entering ValidateEmailFormat")

	if email == "" || !pattern.MatchString(email) {
		log.Warn().Str("email", email).Msg("Invalid email format")
		return apperrors.InvalidEmailFormatErrorDefault
	}

	log.Info().Str("email", email).Msg("Email format is valid")
	return nil
}

// CheckUserExists checks if a user with the specified email exists in the repository.
//
// Parameters:
//   - userRepo: The repository to access user data.
//   - email: The email of the user to check for existence.
//
// Returns:
//   - error: Returns nil if the user does not exist; if the user exists, returns an EmailAlreadyExistsError.
//     If there is an error querying the user repository, it returns that error.
func CheckUserExists(userRepo api.UserRepository, email string) error {
	log.Debug().Str("email", email).Msg("Entering CheckUserExists")

	user, err := userRepo.GetUserByEmail(email)
	if err == nil && user != nil {
		log.Warn().Str("email", email).Msg("Email already exists")
		return apperrors.EmailAlreadyExistsErrorDefault
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Info().Str("email", email).Msg("Email not found, user does not exist")
		return nil
	}

	log.Error().Err(err).Str("email", email).Msg("Error occurred during user existence check")
	return err
}

// ValidatePasswords checks if the new password and its confirmation are valid and match each other.
//
// Parameters:
//   - password: The new password to validate.
//   - passwordConfirmation: The password confirmation to validate.
//
// Returns:
//   - error: Returns nil if the passwords are valid and match; otherwise, returns an InvalidPasswordError.
func ValidatePasswords(password, passwordConfirmation string) error {
	log.Debug().Msg("Entering ValidatePasswords")

	if password == "" || passwordConfirmation == "" {
		log.Warn().Msg("Password or confirmation is empty")
		return apperrors.NewInvalidPasswordError("Passwords can't have empty value")
	}

	if password != passwordConfirmation {
		log.Warn().Msg("Passwords do not match")
		return apperrors.NewInvalidPasswordError("Passwords do not match")
	}

	log.Info().Msg("Passwords are valid and match")
	return nil
}

// ValidatePasswordUpdate validates the password update process for a user.
//
// Parameters:
//   - userUpdateDTO: The data transfer object containing user update information.
//   - user: The user whose password is being updated.
//   - passwordEncoder: The password encoder to validate the current password.
//
// Returns:
//   - error: Returns nil if the password update validation is successful; otherwise, returns an InvalidPasswordError.
func ValidatePasswordUpdate(userUpdateDTO *dto.UserUpdateDTO, user *models.User, passwordEncoder api.PasswordEncoder) error {
	log.Debug().Int64("userID", user.UserID).Msg("Entering ValidatePasswordUpdate")

	matches, err := passwordEncoder.Matches(userUpdateDTO.CurrentPassword, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Int64("userID", user.UserID).Msg("Password matching failed")
		return apperrors.NewInvalidPasswordError("Matcher Failed matching process")
	}

	if !matches {
		log.Warn().Int64("userID", user.UserID).Msg("Current password is incorrect")
		return apperrors.NewInvalidPasswordError("Current password is incorrect")
	}

	if userUpdateDTO.NewPassword == userUpdateDTO.CurrentPassword {
		log.Warn().Int64("userID", user.UserID).Msg("New password is same as current password")
		return apperrors.NewInvalidPasswordError("New password cannot be the same as the current password")
	}

	if err := ValidatePasswords(userUpdateDTO.NewPassword, userUpdateDTO.NewPasswordConfirmation); err != nil {
		log.Error().Err(err).Int64("userID", user.UserID).Msg("Password validation failed")
		return err
	}

	log.Info().Int64("userID", user.UserID).Msg("Password update validation successful")
	return nil
}

// ValidateAuthenticatedUserID checks if the provided user ID matches the authenticated user ID.
//
// Parameters:
//   - userIdFromDto: The user ID from the data transfer object.
//   - userID: The authenticated user's ID.
//
// Returns:
//   - error: Returns nil if the IDs match; otherwise, returns an AccessDeniedError.
func ValidateAuthenticatedUserID(userIdFromDto int64, userID int64) error {
	log.Debug().Int64("userIdFromDto", userIdFromDto).Int64("userID", userID).Msg("Entering ValidateAuthenticatedUserID")

	if userIdFromDto == 0 || userID == 0 || userIdFromDto != userID {
		log.Warn().Int64("userIdFromDto", userIdFromDto).Int64("userID", userID).Msg("User ID mismatch or invalid")
		return apperrors.NewAccessDeniedError("User is not authorized to perform this action")
	}

	log.Info().Int64("userIdFromDto", userIdFromDto).Int64("userID", userID).Msg("Authenticated user ID validation successful")
	return nil
}

// ValidateUserLoginDto validates the user login data transfer object.
//
// Parameters:
//   - userLoginDto: The data transfer object containing user login information.
//
// Returns:
//   - error: Returns nil if the DTO is valid; otherwise, returns an IllegalArgumentError.
func ValidateUserLoginDto(userLoginDto *dto.UserLoginDTO) error {
	log.Debug().Msg("Entering ValidateUserLoginDto")

	if userLoginDto == nil {
		log.Warn().Msg("UserLoginDTO is nil")
		return apperrors.NewIllegalArgumentError("UserLoginDTO is nil")
	}
	if userLoginDto.Email == "" {
		log.Warn().Msg("UserLoginDTO email is empty")
		return apperrors.NewIllegalArgumentError("UserLoginDTO email is nil or empty")
	}
	if userLoginDto.Password == "" {
		log.Warn().Msg("UserLoginDTO password is empty")
		return apperrors.NewIllegalArgumentError("UserLoginDTO password is nil or empty")
	}

	log.Info().Str("email", userLoginDto.Email).Msg("UserLoginDTO validation successful")
	return nil
}

// ValidateUserCreationDTO validates the user creation data transfer object.
//
// Parameters:
//   - userCreationDTO: The data transfer object containing user creation information.
//
// Returns:
//   - error: Returns nil if the DTO is valid; otherwise, returns an IllegalArgumentError or InvalidPasswordError.
func ValidateUserCreationDTO(userCreationDTO *dto.UserCreationDTO) error {
	log.Debug().Msg("Entering ValidateUserCreationDTO")

	if userCreationDTO == nil {
		log.Warn().Msg("UserCreationDTO is nil")
		return apperrors.NewIllegalArgumentError("UserCreationDTO is nil")
	}
	if userCreationDTO.Email == "" {
		log.Warn().Msg("UserCreationDTO email is empty")
		return apperrors.NewIllegalArgumentError("UserCreationDTO email is nil or empty")
	}
	if userCreationDTO.Password == "" {
		log.Warn().Msg("UserCreationDTO password is empty")
		return apperrors.NewIllegalArgumentError("UserCreationDTO password is nil or empty")
	}
	if userCreationDTO.PasswordConfirmation == "" {
		log.Warn().Msg("UserCreationDTO password confirmation is empty")
		return apperrors.NewIllegalArgumentError("UserCreationDTO password confirmation is nil or empty")
	}

	if err := ValidatePasswords(userCreationDTO.Password, userCreationDTO.PasswordConfirmation); err != nil {
		log.Error().Err(err).Msg("Password validation failed for UserCreationDTO")
		return err
	}

	log.Info().Str("email", userCreationDTO.Email).Msg("UserCreationDTO validation successful")
	return nil
}

// ValidateUserUpdateDTO validates the user update data transfer object.
//
// Parameters:
//   - userUpdateDTO: The data transfer object containing user update information.
//
// Returns:
//   - error: Returns nil if the DTO is valid; otherwise, returns an IllegalArgumentError.
func ValidateUserUpdateDTO(userUpdateDTO *dto.UserUpdateDTO) error {
	log.Debug().Msg("Entering ValidateUserUpdateDTO")

	if userUpdateDTO == nil {
		log.Warn().Msg("UserUpdateDTO is nil")
		return apperrors.NewIllegalArgumentError("UserUpdateDTO is nil")
	}
	if userUpdateDTO.CurrentPassword == "" {
		log.Warn().Msg("UserUpdateDTO currentPassword is empty")
		return apperrors.NewIllegalArgumentError("UserUpdateDTO currentPassword is nil or empty")
	}
	if userUpdateDTO.NewPassword == "" {
		log.Warn().Msg("UserUpdateDTO newPassword is empty")
		return apperrors.NewIllegalArgumentError("UserUpdateDTO newPassword is nil or empty")
	}
	if userUpdateDTO.NewPasswordConfirmation == "" {
		log.Warn().Msg("UserUpdateDTO newPasswordConfirmation is empty")
		return apperrors.NewIllegalArgumentError("UserUpdateDTO newPasswordConfirmation is nil or empty")
	}
	if userUpdateDTO.NewPasswordConfirmation != userUpdateDTO.NewPassword {
		log.Warn().Msg("UserUpdateDTO newPasswordConfirmation is not equal to newPassword")
		return apperrors.NewIllegalArgumentError("Passwords do not match")
	}

	log.Info().Msg("UserUpdateDTO validation successful")
	return nil
}

// ValidateUserDeletionDTO validates the user deletion data transfer object.
//
// Parameters:
//   - userDeletionDTO: The data transfer object containing user deletion information.
//
// Returns:
//   - error: Returns nil if the DTO is valid; otherwise, returns an IllegalArgumentError.
func ValidateUserDeletionDTO(userDeletionDTO *dto.UserDeletionDTO) error {
	log.Debug().Msg("Entering ValidateUserDeletionDTO")

	if userDeletionDTO == nil {
		log.Warn().Msg("UserDeletionDTO is nil")
		return apperrors.NewIllegalArgumentError("UserDeletionDTO is nil")
	}
	if userDeletionDTO.Email == "" {
		log.Warn().Msg("UserDeletionDTO email is empty")
		return apperrors.NewIllegalArgumentError("UserDeletionDTO email is nil or empty")
	}
	if userDeletionDTO.CurrentPassword == "" {
		log.Warn().Msg("UserDeletionDTO currentPassword is empty")
		return apperrors.NewIllegalArgumentError("UserDeletionDTO password is nil or empty")
	}

	log.Info().Str("email", userDeletionDTO.Email).Msg("UserDeletionDTO validation successful")
	return nil
}
