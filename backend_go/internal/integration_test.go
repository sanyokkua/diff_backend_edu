package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go_backend/internal/config"
	"go_backend/internal/controllers"
	"go_backend/internal/dto"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const (
	dbName     = "testdb"
	dbUser     = "user"
	dbPassword = "password"
	jwtSecret  = "fdbjchewbdfuywefywtgveytweuywgrtygfvytr"
)

// TestContext holds the shared state for test cases.
type TestContext struct {
	User1Token   string
	User2Token   string
	User1ID      int64
	User2ID      int64
	User1Task1Id int64
	User1Task2Id int64
	User1Task3Id int64
}

func initializeTestEnvironment() (testcontainers.Container, *gin.Engine, error) {
	container, dbPort, err := getDBForTests()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set up test DB: %w", err)
	}

	setupEnvironmentVariables(dbPort)
	router := setupRouter()

	return container, router, nil
}

func getDBForTests() (testcontainers.Container, nat.Port, error) {
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
		return nil, "", fmt.Errorf("failed to start container: %w", err)
	}

	dbPort, err := postgresContainer.MappedPort(c, "5432")
	if err != nil {
		return nil, "", fmt.Errorf("failed to map port: %w", err)
	}

	return postgresContainer, dbPort, nil
}

//goland:noinspection GoUnhandledErrorResult
func setupEnvironmentVariables(dbPort nat.Port) {
	os.Setenv(config.DbHost, "localhost")
	os.Setenv(config.DbUser, dbUser)
	os.Setenv(config.DbPassword, dbPassword)
	os.Setenv(config.DbName, dbName)
	os.Setenv(config.DbPort, dbPort.Port())
	os.Setenv(config.JwtSecret, jwtSecret)
}

func setupRouter() *gin.Engine {
	dependencyContainer, dependencyCreationErr := config.NewDependencyContainer()
	if dependencyCreationErr != nil {
		log.Fatalf("Failed to create dependency container: %v", dependencyCreationErr)
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

func terminateContainer(container testcontainers.Container) {
	if err := container.Terminate(context.Background()); err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}
}

// Helper function to extract data from JSON response
//
//goland:noinspection ALL
func extractToken(response string) string {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return res.Data.(map[string]interface{})["jwtToken"].(string)
}

//goland:noinspection ALL
func extractUserId(response string) int64 {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return int64(res.Data.(map[string]interface{})["userId"].(float64))
}

//goland:noinspection ALL
func extractTaskId(response string) int64 {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return int64(res.Data.(map[string]interface{})["taskId"].(float64))
}

//goland:noinspection ALL
func extractTaskName(response string) string {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return res.Data.(map[string]interface{})["name"].(string)
}

//goland:noinspection ALL
func extractTaskDescription(response string) string {
	var res dto.ResponseDto[any]
	json.Unmarshal([]byte(response), &res)
	return res.Data.(map[string]interface{})["description"].(string)
}

//goland:noinspection ALL
func extractTasks(response string) []dto.TaskDTO {
	var res dto.ResponseDto[[]dto.TaskDTO]
	json.Unmarshal([]byte(response), &res)
	return res.Data
}

func makeRequest(method, url string, body interface{}, authToken string, router *gin.Engine) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestAppIntegration(t *testing.T) {
	container, router, err := initializeTestEnvironment()
	if err != nil {
		log.Println("Failed to set up test environment")
		t.FailNow()
	}
	defer terminateContainer(container)

	ctx := &TestContext{}

	tests := []struct {
		name           string
		method         string
		url            func(*TestContext) string
		body           interface{}
		expectedStatus int
		authToken      func(*TestContext) string
		testCallback   func(*httptest.ResponseRecorder, *TestContext)
	}{
		{
			name:   "Register User 1",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/auth/register"
			},
			body:           dto.UserCreationDTO{Email: "testUser1@email.com", Password: "testUser1Password", PasswordConfirmation: "testUser1Password"},
			expectedStatus: http.StatusCreated,
			authToken:      func(ctx *TestContext) string { return "" },
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				ctx.User1ID = extractUserId(w.Body.String())
				assert.NotZero(t, ctx.User1ID)
			},
		},
		{
			name:   "Register User 2",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/auth/register"
			},
			body:           dto.UserCreationDTO{Email: "testUser2@email.com", Password: "testUser2Password", PasswordConfirmation: "testUser2Password"},
			expectedStatus: http.StatusCreated,
			authToken:      func(ctx *TestContext) string { return "" },
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				ctx.User2ID = extractUserId(w.Body.String())
				assert.NotZero(t, ctx.User2ID)
			},
		},
		{
			name:   "Login User 1",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/auth/login"
			},
			body:           dto.UserLoginDTO{Email: "testUser1@email.com", Password: "testUser1Password"},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ""
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				token := extractToken(w.Body.String())
				ctx.User1Token = token
				assert.NotZero(t, ctx.User1Token)
			},
		},
		{
			name:   "Login User 2",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/auth/login"
			},
			body:           dto.UserLoginDTO{Email: "testUser2@email.com", Password: "testUser2Password"},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ""
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				token := extractToken(w.Body.String())
				ctx.User2Token = token
				assert.NotZero(t, ctx.User2Token)
			},
		},
		{
			name:   "User 1 changed password",
			method: http.MethodPut,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/password"
			},
			body:           dto.UserUpdateDTO{CurrentPassword: "testUser1Password", NewPassword: "newPassword1", NewPasswordConfirmation: "newPassword1"},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "Login User 1 with new password",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/auth/login"
			},
			body:           dto.UserLoginDTO{Email: "testUser1@email.com", Password: "newPassword1"},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ""
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				token := extractToken(w.Body.String())
				ctx.User1Token = token
				assert.NotZero(t, ctx.User1Token)
			},
		},
		{
			name:   "User 2 create task",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User2ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task 1", Description: "Task 1 Description"},
			expectedStatus: http.StatusCreated,
			authToken: func(ctx *TestContext) string {
				return ctx.User2Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				assert.NotZero(t, taskId)
			},
		},
		{
			name:   "User 2 delete account",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User2ID, 10) + "/delete"
			},
			body:           dto.UserDeletionDTO{Email: "testUser2@email.com", CurrentPassword: "testUser2Password"},
			expectedStatus: http.StatusNoContent,
			authToken: func(ctx *TestContext) string {
				return ctx.User2Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 create task",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task 1", Description: "Task 1 Description"},
			expectedStatus: http.StatusCreated,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				ctx.User1Task1Id = taskId
				assert.NotZero(t, taskId)
			},
		},
		{
			name:   "User 1 update task",
			method: http.MethodPut,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task1Id, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskUpdateDTO{Name: "Updated Task Name", Description: "Updated Task Description"},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				taskName := extractTaskName(w.Body.String())
				taskDescription := extractTaskDescription(w.Body.String())
				assert.NotZero(t, taskId)
				assert.NotZero(t, taskName)
				assert.NotZero(t, taskDescription)
				assert.Equal(t, "Updated Task Name", taskName)
				assert.Equal(t, "Updated Task Description", taskDescription)
			},
		},
		{
			name:   "User 1 delete task",
			method: http.MethodDelete,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task1Id, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusNoContent,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 create task new 1",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task 1", Description: "Task 1 Description"},
			expectedStatus: http.StatusCreated,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				ctx.User1Task1Id = taskId
				assert.NotZero(t, taskId)
			},
		},
		{
			name:   "User 1 create task new 2",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task 2", Description: "Task 2 Description"},
			expectedStatus: http.StatusCreated,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				ctx.User1Task2Id = taskId
				assert.NotZero(t, taskId)
			},
		},
		{
			name:   "User 1 create task new 3",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task 3", Description: "Task 3 Description"},
			expectedStatus: http.StatusCreated,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				taskId := extractTaskId(w.Body.String())
				ctx.User1Task3Id = taskId
				assert.NotZero(t, taskId)
			},
		},
		{
			name:   "User 1 get all tasks (len 3)",
			method: http.MethodGet,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				body := w.Body.String()
				tasks := extractTasks(body)
				assert.NotZero(t, body)
				assert.Len(t, tasks, 3)
			},
		},
		{
			name:   "User 1 delete task 1",
			method: http.MethodDelete,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task1Id, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusNoContent,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 get all tasks (len 2)",
			method: http.MethodGet,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				body := w.Body.String()
				tasks := extractTasks(body)
				assert.NotZero(t, body)
				assert.Len(t, tasks, 2)
			},
		},
		{
			name:   "User 1 delete non existing task",
			method: http.MethodDelete,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(9999, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusNotFound,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 create task with empty fields",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "", Description: ""},
			expectedStatus: http.StatusBadRequest,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 create task without authorization token",
			method: http.MethodPost,
			url: func(ctx *TestContext) string {
				return "/api/v1/users/" + strconv.FormatInt(ctx.User1ID, 10) + "/tasks/"
			},
			body:           dto.TaskCreationDTO{Name: "Task", Description: "Description"},
			expectedStatus: http.StatusUnauthorized,
			authToken: func(ctx *TestContext) string {
				return ""
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 2 delete task of user 1",
			method: http.MethodDelete,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task1Id, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusUnauthorized,
			authToken: func(ctx *TestContext) string {
				return ctx.User2Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 update task with incorrect data",
			method: http.MethodPut,
			url: func(ctx *TestContext) string {
				userIdStr := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task1Id, 10)
				return "/api/v1/users/" + userIdStr + "/tasks/" + taskIdStr
			},
			body:           dto.TaskUpdateDTO{Name: "", Description: ""},
			expectedStatus: http.StatusBadRequest,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
			},
		},
		{
			name:   "User 1 get single task",
			method: http.MethodGet,
			url: func(ctx *TestContext) string {
				userId := strconv.FormatInt(ctx.User1ID, 10)
				taskIdStr := strconv.FormatInt(ctx.User1Task2Id, 10)
				return "/api/v1/users/" + userId + "/tasks/" + taskIdStr
			},
			body:           dto.TaskDTO{},
			expectedStatus: http.StatusOK,
			authToken: func(ctx *TestContext) string {
				return ctx.User1Token
			},
			testCallback: func(w *httptest.ResponseRecorder, ctx *TestContext) {
				body := w.Body.String()
				taskName := extractTaskName(body)
				taskDesc := extractTaskDescription(body)
				assert.NotZero(t, body)
				assert.NotZero(t, taskName)
				assert.NotZero(t, taskDesc)
				assert.Equal(t, "Task 2", taskName)
				assert.Equal(t, "Task 2 Description", taskDesc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := makeRequest(tt.method, tt.url(ctx), tt.body, tt.authToken(ctx), router)
			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.testCallback(w, ctx)
		})
	}
}
