package services

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
)

func TestUserService(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordEncoder := new(MockPasswordEncoder)
	userService := NewUserService(userRepo, passwordEncoder)

	t.Run("Create_Success", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userCreationDTO := &dto.UserCreationDTO{
			Email:                "test@example.com",
			Password:             "password123",
			PasswordConfirmation: "password123",
		}

		encodedPassword := "encodedPassword"
		newUser := models.NewFullUser(0, userCreationDTO.Email, encodedPassword, nil)
		createdUser := &models.User{
			UserID:       1,
			Email:        userCreationDTO.Email,
			PasswordHash: encodedPassword,
		}

		var user *models.User = nil

		passwordEncoder.On("Encode", userCreationDTO.Password).Return(encodedPassword, nil)
		userRepo.On("CreateUser", newUser).Return(createdUser, nil)
		userRepo.On("GetUserByEmail", userCreationDTO.Email).Return(user, gorm.ErrRecordNotFound)

		userDTO, err := userService.Create(userCreationDTO)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, createdUser.UserID, userDTO.UserID)
		assert.Equal(t, createdUser.Email, userDTO.Email)

		userRepo.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("Create_InvalidEmail", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userCreationDTO := &dto.UserCreationDTO{
			Email:                "invalid-email",
			Password:             "password123",
			PasswordConfirmation: "password123",
		}

		_, err := userService.Create(userCreationDTO)

		assert.Error(t, err)
		assert.IsType(t, apperrors.InvalidEmailFormatError{}, err)
	})

	t.Run("UpdatePassword_Success", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userID := int64(1)
		userUpdateDTO := &dto.UserUpdateDTO{
			CurrentPassword:         "currentPassword",
			NewPassword:             "newPassword123",
			NewPasswordConfirmation: "newPassword123",
		}

		storedUser := &models.User{
			UserID:       userID,
			Email:        "test@example.com",
			PasswordHash: "hashedCurrentPassword",
		}

		passwordEncoder.On("Matches", userUpdateDTO.CurrentPassword, storedUser.PasswordHash).Return(true, nil)
		passwordEncoder.On("Encode", userUpdateDTO.NewPassword).Return("encodedNewPassword", nil)
		userRepo.On("GetUserByID", userID).Return(storedUser, nil)
		userRepo.On("UpdateUser", mock.Anything).Return(storedUser, nil)

		userDTO, err := userService.UpdatePassword(userID, userUpdateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, storedUser.UserID, userDTO.UserID)
		assert.Equal(t, storedUser.Email, userDTO.Email)

		userRepo.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("UpdatePassword_InvalidCurrentPassword", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userID := int64(1)
		userUpdateDTO := &dto.UserUpdateDTO{
			CurrentPassword:         "wrongCurrentPassword",
			NewPassword:             "newPassword123",
			NewPasswordConfirmation: "newPassword123",
		}

		storedUser := &models.User{
			UserID:       userID,
			Email:        "test@example.com",
			PasswordHash: "hashedCurrentPassword",
		}

		passwordEncoder.On("Matches", userUpdateDTO.CurrentPassword, storedUser.PasswordHash).Return(false, nil)
		userRepo.On("GetUserByID", userID).Return(storedUser, nil)

		userDTO, err := userService.UpdatePassword(userID, userUpdateDTO)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.IsType(t, apperrors.InvalidPasswordError{}, err)

		userRepo.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("Delete_Success", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userID := int64(1)
		userDeletionDTO := &dto.UserDeletionDTO{
			Email:           "test@example.com",
			CurrentPassword: "currentPassword",
		}

		storedUser := &models.User{
			UserID:       userID,
			Email:        "test@example.com",
			PasswordHash: "hashedCurrentPassword",
		}

		passwordEncoder.On("Matches", userDeletionDTO.CurrentPassword, storedUser.PasswordHash).Return(true, nil)
		userRepo.On("GetUserByID", userID).Return(storedUser, nil)
		userRepo.On("DeleteUser", userID).Return(nil)

		err := userService.Delete(userID, userDeletionDTO)

		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("Delete_InvalidCurrentPassword", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userID := int64(1)
		userDeletionDTO := &dto.UserDeletionDTO{
			Email:           "test@example.com",
			CurrentPassword: "wrongCurrentPassword",
		}

		storedUser := &models.User{
			UserID:       userID,
			Email:        "test@example.com",
			PasswordHash: "hashedCurrentPassword",
		}

		passwordEncoder.On("Matches", userDeletionDTO.CurrentPassword, storedUser.PasswordHash).Return(false, nil)
		userRepo.On("GetUserByID", userID).Return(storedUser, nil)

		err := userService.Delete(userID, userDeletionDTO)

		assert.Error(t, err)
		assert.IsType(t, apperrors.InvalidPasswordError{}, err)

		userRepo.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("Delete_UserNotFound", func(t *testing.T) {
		t.Cleanup(func() {
			userRepo = new(MockUserRepository)
			passwordEncoder = new(MockPasswordEncoder)
			userService = NewUserService(userRepo, passwordEncoder)
		})
		userID := int64(1)
		userDeletionDTO := &dto.UserDeletionDTO{
			Email:           "test@example.com",
			CurrentPassword: "currentPassword",
		}

		var userNil *models.User = nil
		userRepo.On("GetUserByID", userID).Return(userNil, errors.New("user not found"))

		err := userService.Delete(userID, userDeletionDTO)

		assert.Error(t, err)

		userRepo.AssertExpectations(t)
	})
}
