package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
	"testing"
)

func TestAuthenticationService(t *testing.T) {
	userService := new(MockUserService)
	userRepository := new(MockUserRepository)
	jwtService := new(MockJwtService)
	passwordEncoder := new(MockPasswordEncoder)

	authService := NewAuthenticationService(userService, userRepository, jwtService, passwordEncoder)

	t.Run("LoginUser_Success", func(t *testing.T) {
		loginDTO := &dto.UserLoginDTO{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUser := &models.User{
			UserID:       1,
			Email:        loginDTO.Email,
			PasswordHash: "hashedPassword",
		}

		userRepository.On("GetUserByEmail", loginDTO.Email).Return(mockUser, nil)
		passwordEncoder.On("Matches", loginDTO.Password, mockUser.PasswordHash).Return(true, nil)
		jwtService.On("GenerateJwtToken", mockUser.Email).Return("jwtToken", nil)

		userDTO, err := authService.LoginUser(loginDTO)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, mockUser.UserID, userDTO.UserID)
		assert.Equal(t, mockUser.Email, userDTO.Email)
		assert.Equal(t, "jwtToken", userDTO.JwtToken)

		userRepository.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("LoginUser_InvalidEmail", func(t *testing.T) {
		loginDTO := &dto.UserLoginDTO{
			Email:    "invalid@example.com",
			Password: "password123",
		}
		var userNil *models.User = nil

		userRepository.On("GetUserByEmail", loginDTO.Email).Return(userNil, errors.New("user not found"))

		userDTO, err := authService.LoginUser(loginDTO)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.IsType(t, apperrors.IllegalArgumentError{}, err)

		userRepository.AssertExpectations(t)
	})

	t.Run("LoginUser_InvalidPassword", func(t *testing.T) {
		loginDTO := &dto.UserLoginDTO{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockUser := &models.User{
			UserID:       1,
			Email:        loginDTO.Email,
			PasswordHash: "hashedPassword",
		}

		userRepository.On("GetUserByEmail", loginDTO.Email).Return(mockUser, nil)
		passwordEncoder.On("Matches", loginDTO.Password, mockUser.PasswordHash).Return(false, nil)

		userDTO, err := authService.LoginUser(loginDTO)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.IsType(t, apperrors.InvalidPasswordError{}, err)

		userRepository.AssertExpectations(t)
		passwordEncoder.AssertExpectations(t)
	})

	t.Run("RegisterUser_Success", func(t *testing.T) {
		userCreationDTO := &dto.UserCreationDTO{
			Email:                "newuser@example.com",
			Password:             "newpassword123",
			PasswordConfirmation: "newpassword123",
		}

		mockUserDTO := &dto.UserDTO{
			UserID:   2,
			Email:    userCreationDTO.Email,
			JwtToken: "newJwtToken",
		}

		userService.On("Create", userCreationDTO).Return(mockUserDTO, nil)
		jwtService.On("GenerateJwtToken", mockUserDTO.Email).Return("newJwtToken", nil)

		userDTO, err := authService.RegisterUser(userCreationDTO)

		assert.NoError(t, err)
		assert.NotNil(t, userDTO)
		assert.Equal(t, mockUserDTO.UserID, userDTO.UserID)
		assert.Equal(t, mockUserDTO.Email, userDTO.Email)
		assert.Equal(t, "newJwtToken", userDTO.JwtToken)

		userService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("RegisterUser_InvalidDetails", func(t *testing.T) {
		userCreationDTO := &dto.UserCreationDTO{
			Email:                "invalidemail",
			Password:             "short",
			PasswordConfirmation: "short",
		}

		var userNil *dto.UserDTO = nil
		userService.On("Create", userCreationDTO).Return(userNil, errors.New("invalid user details"))

		userDTO, err := authService.RegisterUser(userCreationDTO)

		assert.Error(t, err)
		assert.Nil(t, userDTO)
		assert.IsType(t, apperrors.IllegalArgumentError{}, err)

		userService.AssertExpectations(t)
	})
}
