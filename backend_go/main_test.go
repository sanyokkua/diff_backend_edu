package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/go-connections/nat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go_backend/config"
	"go_backend/controllers"
	"go_backend/dto"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const dbName = "testdb"
const dbUser = "user"
const dbPassword = "password"

func GetDBForTests() (testcontainers.Container, nat.Port, error) {
	c := context.Background()

	postgresContainer, err := postgres.Run(c,
		"docker.io/postgres",
		postgres.WithInitScripts(filepath.Join("./testdata", "init.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, "", err
	}

	dbPort, err := postgresContainer.MappedPort(c, "5432")
	if err != nil {
		return nil, "", err
	}

	return postgresContainer, dbPort, nil
}

// Test structure
type AppApplicationTests struct {
	DB      *gorm.DB
	Router  *gin.Engine
	User1ID uint
	User2ID uint
}

func SetupRouter() *gin.Engine {
	dependencyContainer, dependencyCreationErr := config.NewDependencyContainer()
	if dependencyCreationErr != nil {
		panic(dependencyContainer)
	}
	authController := controllers.NewAuthController(dependencyContainer.AuthenticationService)
	userController := controllers.NewUserController(dependencyContainer.UserService)
	taskController := controllers.NewTaskController(dependencyContainer.TaskService)
	router := gin.Default()

	router.Use(cors.New(*dependencyContainer.CorsConfig))
	controllers.RegisterAuthRoutes(router, authController)
	controllers.RegisterUserRoutes(router, userController, dependencyContainer.JwtMiddleware)
	controllers.RegisterTaskRoutes(router, taskController, dependencyContainer.JwtMiddleware)

	return router
}

// Helper function to extract data from JSON response
func extractToken(response string) string {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return res.Data.(map[string]interface{})["jwtToken"].(string)
}

func extractUserId(response string) int64 {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return int64(res.Data.(map[string]interface{})["userId"].(float64))
}

func TestAppIntegration(t *testing.T) {
	container, dbPort, err := GetDBForTests()
	if err != nil {
		log.Println("Failed to set up test DB")
		t.FailNow()
	}
	defer func() {
		if err := container.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	os.Setenv(config.DbHost, "localhost")
	os.Setenv(config.DbUser, dbUser)
	os.Setenv(config.DbPassword, dbPassword)
	os.Setenv(config.DbName, dbName)
	os.Setenv(config.DbPort, dbPort.Port())
	os.Setenv(config.JwtSecret, "fdbjchewbdfuywefywtgveytweuywgrtygfvytr")

	// Initialize the router
	router := SetupRouter()

	user1Token := ""
	user2Token := ""
	var user1Id int64
	var user2Id int64
	//var user1Task1Id int64
	//var user1Task2Id int64
	//var user1Task3Id int64
	// Store state across tests
	tests := []struct {
		name           string
		method         string
		url            string
		body           interface{}
		expectedStatus int
		expectedUserId *uint
		testCallback   func(*httptest.ResponseRecorder)
	}{
		{
			name:           "Register User 1",
			method:         http.MethodPost,
			url:            "/api/v1/auth/register",
			body:           dto.UserCreationDTO{Email: "testUser1@email.com", Password: "testUser1Password", PasswordConfirmation: "testUser1Password"},
			expectedStatus: http.StatusCreated,
			expectedUserId: nil, // Will be populated after the request
			testCallback: func(w *httptest.ResponseRecorder) {
				user1Id = extractUserId(w.Body.String())
				assert.NotZero(t, user1Id)
			},
		},
		{
			name:           "Register User 2",
			method:         http.MethodPost,
			url:            "/api/v1/auth/register",
			body:           dto.UserCreationDTO{Email: "testUser2@email.com", Password: "testUser2Password", PasswordConfirmation: "testUser2Password"},
			expectedStatus: http.StatusCreated,
			expectedUserId: nil, // Will be populated after the request
			testCallback: func(w *httptest.ResponseRecorder) {
				user2Id = extractUserId(w.Body.String())
				assert.NotZero(t, user2Id)
			},
		},
		{
			name:           "Login User 1",
			method:         http.MethodPost,
			url:            "/api/v1/auth/login",
			body:           dto.UserLoginDTO{Email: "testUser1@email.com", Password: "testUser1Password"},
			expectedStatus: http.StatusOK,
			expectedUserId: nil, // Will be populated after the request
			testCallback: func(w *httptest.ResponseRecorder) {
				user1Token = extractToken(w.Body.String())
				assert.NotZero(t, user1Token)
			},
		},
		{
			name:           "Login User 2",
			method:         http.MethodPost,
			url:            "/api/v1/auth/login",
			body:           dto.UserLoginDTO{Email: "testUser2@email.com", Password: "testUser2Password"},
			expectedStatus: http.StatusOK,
			expectedUserId: nil, // Will be populated after the request
			testCallback: func(w *httptest.ResponseRecorder) {
				user2Token = extractToken(w.Body.String())
				assert.NotZero(t, user2Token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.body != nil {
				jsonData, _ := json.Marshal(tt.body)
				reqBody = bytes.NewBuffer(jsonData)
			}

			req, _ := http.NewRequest(tt.method, tt.url, reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.testCallback(w)
		})
	}
}
