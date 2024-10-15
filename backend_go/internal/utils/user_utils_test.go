package utils

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for ValidateEmailFormat
func TestValidateEmailFormat(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectedErr error
	}{
		{"Valid Email", "valid.email@example.com", nil},
		{"Invalid Email", "invalid-email.com", apperrors.InvalidEmailFormatErrorDefault},
		{"Empty Email", "", apperrors.InvalidEmailFormatErrorDefault},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailFormat(tt.email)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Mocking the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int64) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

// Test for CheckUserExists
func TestCheckUserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	tests := []struct {
		name         string
		email        string
		userReturned *models.User
		errReturned  error
		expectedErr  error
	}{
		{"User Exists", "existing@example.com", &models.User{}, nil, apperrors.EmailAlreadyExistsErrorDefault},
		{"User Does Not Exist", "nonexistent@example.com", nil, gorm.ErrRecordNotFound, nil},
		{"Database Error", "error@example.com", nil, errors.New("db error"), errors.New("db error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUserByEmail", tt.email).Return(tt.userReturned, tt.errReturned)
			err := CheckUserExists(mockRepo, tt.email)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Test for ValidatePasswords
func TestValidatePasswords(t *testing.T) {
	tests := []struct {
		name                 string
		password             string
		passwordConfirmation string
		expectedErr          error
	}{
		{"Matching Passwords", "password123", "password123", nil},
		{"Non-Matching Passwords", "password123", "differentpassword", apperrors.NewInvalidPasswordError("Passwords do not match")},
		{"Empty Password", "", "password123", apperrors.NewInvalidPasswordError("Passwords can't have empty value")},
		{"Empty Confirmation", "password123", "", apperrors.NewInvalidPasswordError("Passwords can't have empty value")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswords(tt.password, tt.passwordConfirmation)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Mocking the PasswordEncoder interface
type MockPasswordEncoder struct {
	mock.Mock
}

func (m *MockPasswordEncoder) Encode(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordEncoder) Matches(rawPassword, encodedPassword string) (bool, error) {
	args := m.Called(rawPassword, encodedPassword)
	return args.Bool(0), args.Error(1)
}

// Test for ValidatePasswordUpdate
func TestValidatePasswordUpdate(t *testing.T) {
	mockPasswordEncoder := new(MockPasswordEncoder)
	user := &models.User{UserID: 1, PasswordHash: "encoded_password"}

	tests := []struct {
		name        string
		dto         *dto.UserUpdateDTO
		matches     bool
		matchError  error
		expectedErr error
	}{
		{"Valid Password Update", &dto.UserUpdateDTO{CurrentPassword: "oldpassword", NewPassword: "newpassword", NewPasswordConfirmation: "newpassword"}, true, nil, nil},
		{"Password Mismatch", &dto.UserUpdateDTO{CurrentPassword: "oldpassword", NewPassword: "newpassword", NewPasswordConfirmation: "wrongconfirmation"}, true, nil, apperrors.NewInvalidPasswordError("Passwords do not match")},
		{"Incorrect Current Password", &dto.UserUpdateDTO{CurrentPassword: "wrongpassword", NewPassword: "newpassword", NewPasswordConfirmation: "newpassword"}, false, nil, apperrors.NewInvalidPasswordError("Current password is incorrect")},
		{"New Password Same as Old", &dto.UserUpdateDTO{CurrentPassword: "oldpassword", NewPassword: "oldpassword", NewPasswordConfirmation: "oldpassword"}, true, nil, apperrors.NewInvalidPasswordError("New password cannot be the same as the current password")},
		{"Password Encoding Error", &dto.UserUpdateDTO{CurrentPassword: "oldpassword", NewPassword: "newpassword2", NewPasswordConfirmation: "newpassword"}, false, errors.New("encoding error"), apperrors.NewInvalidPasswordError("Passwords do not match")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPasswordEncoder.On("Matches", tt.dto.CurrentPassword, user.PasswordHash).Return(tt.matches, tt.matchError)
			err := ValidatePasswordUpdate(tt.dto, user, mockPasswordEncoder)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Test for ValidateAuthenticatedUserID
func TestValidateAuthenticatedUserID(t *testing.T) {
	tests := []struct {
		name          string
		userIdFromDto int64
		userID        int64
		expectedErr   error
	}{
		{"Authorized User", 1, 1, nil},
		{"Unauthorized User", 1, 2, apperrors.NewAccessDeniedError("User is not authorized to perform this action")},
		{"Empty User ID from DTO", 0, 1, apperrors.NewAccessDeniedError("User is not authorized to perform this action")},
		{"Empty Authenticated User ID", 1, 0, apperrors.NewAccessDeniedError("User is not authorized to perform this action")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthenticatedUserID(tt.userIdFromDto, tt.userID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Test for ValidateUserLoginDto
func TestValidateUserLoginDto(t *testing.T) {
	tests := []struct {
		name          string
		userLoginDto  *dto.UserLoginDTO
		expectedError error
	}{
		{"Valid UserLoginDTO", &dto.UserLoginDTO{Email: "test@example.com", Password: "password123"}, nil},
		{"Nil UserLoginDTO", nil, apperrors.NewIllegalArgumentError("UserLoginDTO is nil")},
		{"Empty Email", &dto.UserLoginDTO{Email: "", Password: "password123"}, apperrors.NewIllegalArgumentError("UserLoginDTO email is nil or empty")},
		{"Empty Password", &dto.UserLoginDTO{Email: "test@example.com", Password: ""}, apperrors.NewIllegalArgumentError("UserLoginDTO password is nil or empty")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserLoginDto(tt.userLoginDto)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestValidateUserCreationDTO(t *testing.T) {
	tests := []struct {
		name            string
		userCreationDto *dto.UserCreationDTO
		expectedError   error
	}{
		{
			name: "Valid UserCreationDTO",
			userCreationDto: &dto.UserCreationDTO{
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			expectedError: nil,
		},
		{
			name:            "Nil UserCreationDTO",
			userCreationDto: nil,
			expectedError:   apperrors.NewIllegalArgumentError("UserCreationDTO is nil"),
		},
		{
			name: "Empty Email",
			userCreationDto: &dto.UserCreationDTO{
				Email:                "",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			expectedError: apperrors.NewIllegalArgumentError("UserCreationDTO email is nil or empty"),
		},
		{
			name: "Empty Password",
			userCreationDto: &dto.UserCreationDTO{
				Email:                "test@example.com",
				Password:             "",
				PasswordConfirmation: "password123",
			},
			expectedError: apperrors.NewIllegalArgumentError("UserCreationDTO password is nil or empty"),
		},
		{
			name: "Password and Confirmation Mismatch",
			userCreationDto: &dto.UserCreationDTO{
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password456",
			},
			expectedError: apperrors.NewInvalidPasswordError("Passwords do not match"),
		},
		{
			name: "Password and Confirmation Mismatch",
			userCreationDto: &dto.UserCreationDTO{
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "",
			},
			expectedError: apperrors.NewInvalidPasswordError("UserCreationDTO password confirmation is nil or empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserCreationDTO(tt.userCreationDto)
			if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestValidateUserUpdateDTO(t *testing.T) {
	tests := []struct {
		name          string
		userUpdateDto *dto.UserUpdateDTO
		expectedError error
	}{
		{
			name: "Valid UserUpdateDTO",
			userUpdateDto: &dto.UserUpdateDTO{
				CurrentPassword:         "oldpassword123",
				NewPassword:             "newpassword123",
				NewPasswordConfirmation: "newpassword123",
			},
			expectedError: nil,
		},
		{
			name:          "Nil UserUpdateDTO",
			userUpdateDto: nil,
			expectedError: apperrors.NewIllegalArgumentError("UserUpdateDTO is nil"),
		},
		{
			name: "Empty Current Password",
			userUpdateDto: &dto.UserUpdateDTO{
				CurrentPassword:         "",
				NewPassword:             "newpassword123",
				NewPasswordConfirmation: "newpassword123",
			},
			expectedError: apperrors.NewIllegalArgumentError("UserUpdateDTO currentPassword is nil or empty"),
		},
		{
			name: "New Password and Confirmation Mismatch",
			userUpdateDto: &dto.UserUpdateDTO{
				CurrentPassword:         "oldpassword123",
				NewPassword:             "newpassword123",
				NewPasswordConfirmation: "newpassword456",
			},
			expectedError: apperrors.NewInvalidPasswordError("Passwords do not match"),
		},
		{
			name: "New Password and Confirmation Mismatch",
			userUpdateDto: &dto.UserUpdateDTO{
				CurrentPassword:         "oldpassword123",
				NewPassword:             "",
				NewPasswordConfirmation: "newpassword456",
			},
			expectedError: apperrors.NewInvalidPasswordError("UserUpdateDTO newPassword is nil or empty"),
		},
		{
			name: "New Password and Confirmation Mismatch",
			userUpdateDto: &dto.UserUpdateDTO{
				CurrentPassword:         "oldpassword123",
				NewPassword:             "newpassword123",
				NewPasswordConfirmation: "",
			},
			expectedError: apperrors.NewInvalidPasswordError("UserUpdateDTO newPasswordConfirmation is nil or empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserUpdateDTO(tt.userUpdateDto)
			if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestValidateUserDeletionDTO(t *testing.T) {
	tests := []struct {
		name            string
		userDeletionDto *dto.UserDeletionDTO
		expectedError   error
	}{
		{
			name: "Valid UserDeletionDTO",
			userDeletionDto: &dto.UserDeletionDTO{
				Email:           "test@example.com",
				CurrentPassword: "password123",
			},
			expectedError: nil,
		},
		{
			name:            "Nil UserDeletionDTO",
			userDeletionDto: nil,
			expectedError:   apperrors.NewIllegalArgumentError("UserDeletionDTO is nil"),
		},
		{
			name: "Empty Email",
			userDeletionDto: &dto.UserDeletionDTO{
				Email:           "",
				CurrentPassword: "password123",
			},
			expectedError: apperrors.NewIllegalArgumentError("UserDeletionDTO email is nil or empty"),
		},
		{
			name: "Empty Current Password",
			userDeletionDto: &dto.UserDeletionDTO{
				Email:           "test@example.com",
				CurrentPassword: "",
			},
			expectedError: apperrors.NewIllegalArgumentError("UserDeletionDTO password is nil or empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserDeletionDTO(tt.userDeletionDto)
			if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
