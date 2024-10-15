package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go_backend/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JWT service and user repository for testing purposes
type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) IsTokenExpired(claims jwt.Claims) bool {
	args := m.Called(claims)
	return args.Bool(0)
}

func (m *MockJwtService) GenerateJwtToken(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) ExtractClaims(tokenString string) (jwt.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) != nil {
		return args.Get(0).(jwt.Claims), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJwtService) ValidateToken(tokenString, subject string) bool {
	args := m.Called(tokenString, subject)
	return args.Bool(0)
}

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
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// Mock Claims to simulate the Claims interface
type MockClaims struct {
	mock.Mock
}

func (m *MockClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*jwt.NumericDate), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*jwt.NumericDate), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClaims) GetNotBefore() (*jwt.NumericDate, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*jwt.NumericDate), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClaims) GetIssuer() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockClaims) GetSubject() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockClaims) GetAudience() (jwt.ClaimStrings, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(jwt.ClaimStrings), args.Error(1)
	}
	return nil, args.Error(1)
}

// Setup a basic Gin context for testing
func setupGinContext(authHeader string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	ctx.Request = req

	return ctx, recorder
}

func TestJWTAuthMiddleware_NoAuthHeader(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("") // No Authorization header

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "authorization header required")
}

func TestJWTAuthMiddleware_MissingBearerToken(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer ") // Missing token after Bearer

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "bearer token required")
}

func TestJWTAuthMiddleware_InvalidClaims(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer invalidtoken")

	// Mock the JWT service to fail claims extraction
	mockJwtService.On("ExtractClaims", "invalidtoken").Return(nil, errors.New("invalid token"))

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid token")
}

func TestJWTAuthMiddleware_InvalidSubject(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer validtoken")

	// Mock valid claims extraction but subject extraction fails
	mockClaims := new(MockClaims)
	mockJwtService.On("ExtractClaims", "validtoken").Return(mockClaims, nil)
	mockClaims.On("GetSubject").Return("", errors.New("invalid subject"))

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid token subject")
}

func TestJWTAuthMiddleware_InvalidOrExpiredToken(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer validtoken")

	// Mock valid claims and subject extraction, but token validation fails
	mockClaims := new(MockClaims)
	mockJwtService.On("ExtractClaims", "validtoken").Return(mockClaims, nil)
	mockClaims.On("GetSubject").Return("user@example.com", nil)
	mockJwtService.On("ValidateToken", "validtoken", "user@example.com").Return(false)

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid or expired token")
}

func TestJWTAuthMiddleware_UserNotFound(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer validtoken")

	// Mock valid claims, subject extraction, and token validation
	mockClaims := new(MockClaims)
	mockJwtService.On("ExtractClaims", "validtoken").Return(mockClaims, nil)
	mockClaims.On("GetSubject").Return("user@example.com", nil)
	mockJwtService.On("ValidateToken", "validtoken", "user@example.com").Return(true)

	// User not found in repository
	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(nil, errors.New("user not found"))

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "user not found")
}

func TestJWTAuthMiddleware_Success(t *testing.T) {
	mockJwtService := new(MockJwtService)
	mockUserRepo := new(MockUserRepository)
	ctx, recorder := setupGinContext("Bearer validtoken")

	// Mock valid claims, subject extraction, and token validation
	mockClaims := new(MockClaims)
	mockJwtService.On("ExtractClaims", "validtoken").Return(mockClaims, nil)
	mockClaims.On("GetSubject").Return("user@example.com", nil)
	mockJwtService.On("ValidateToken", "validtoken", "user@example.com").Return(true)

	// Mock valid user retrieval
	user := &models.User{
		Email: "user@example.com",
	}
	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)

	middleware := JWTAuthMiddleware(mockJwtService, mockUserRepo)
	middleware(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, user, ctx.MustGet("CurrentUser").(*models.User))
}
