package services

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"go_backend/dto"
	"go_backend/models"
)

// MockUserService is a mock implementation of UserService for testing purposes.
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(userCreationDTO *dto.UserCreationDTO) (*dto.UserDTO, error) {
	args := m.Called(userCreationDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdatePassword(userID int64, userUpdateDTO *dto.UserUpdateDTO) (*dto.UserDTO, error) {
	args := m.Called(userID, userUpdateDTO)
	return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func (m *MockUserService) Delete(userID int64, userDeletionDTO *dto.UserDeletionDTO) error {
	args := m.Called(userID, userDeletionDTO)
	return args.Error(0)
}

// MockUserRepository is a mock implementation of UserRepository for testing purposes.
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

// MockJwtService is a mock implementation of JwtService for testing purposes.
type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) ExtractClaims(token string) (jwt.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(jwt.Claims), args.Error(1)
}

func (m *MockJwtService) IsTokenExpired(claims jwt.Claims) bool {
	args := m.Called(claims)
	return args.Bool(0)
}

func (m *MockJwtService) ValidateToken(token, username string) bool {
	args := m.Called(token, username)
	return args.Bool(0)
}

func (m *MockJwtService) GenerateJwtToken(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

// MockPasswordEncoder is a mock implementation of PasswordEncoder for testing purposes.
type MockPasswordEncoder struct {
	mock.Mock
}

func (m *MockPasswordEncoder) Matches(rawPassword, encodedPassword string) (bool, error) {
	args := m.Called(rawPassword, encodedPassword)
	return args.Bool(0), args.Error(1)
}

func (m *MockPasswordEncoder) Encode(rawPassword string) (string, error) {
	args := m.Called(rawPassword)
	return args.String(0), args.Error(1)
}

// Mocked TaskRepository and UserRepository for testing purposes.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindByTaskID(taskID int64) (*models.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUserAndName(user *models.User, name string) (*models.Task, error) {
	args := m.Called(user, name)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUserAndTaskID(user *models.User, taskID int64) (*models.Task, error) {
	args := m.Called(user, taskID)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAllByUser(user *models.User) ([]models.Task, error) {
	args := m.Called(user)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(task *models.Task) error {
	return m.Called(task).Error(0)
}
