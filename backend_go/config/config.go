package config

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_backend/api"
	"go_backend/middlewares"
	"go_backend/repositories"
	"go_backend/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

const (
	DbHost     = "DB_HOST"     // Environment variable key for database host
	DbUser     = "DB_USER"     // Environment variable key for database user
	DbPassword = "DB_PASSWORD" // Environment variable key for database password
	DbName     = "DB_NAME"     // Environment variable key for database name
	DbPort     = "DB_PORT"     // Environment variable key for database port
	JwtSecret  = "JWT_SECRET"  // Environment variable key for JWT secret
)

// getEnvWithDefault retrieves the value of the environment variable with the given key.
// If the variable is not set, it returns the specified default value.
//
// Parameters:
//   - key: the name of the environment variable to retrieve.
//   - defaultValue: the value to return if the environment variable is not set.
//
// Returns:
//   - string: the value of the environment variable or the default value.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetConfig returns a map containing configuration values read from environment variables.
// It uses default values if the environment variables are not set.
func GetConfig() map[string]string {
	return map[string]string{
		DbHost:     getEnvWithDefault(DbHost, "localhost"),
		DbUser:     getEnvWithDefault(DbUser, "development"),
		DbPassword: getEnvWithDefault(DbPassword, "dev_pass"),
		DbName:     getEnvWithDefault(DbName, "backend_diff_db"),
		DbPort:     getEnvWithDefault(DbPort, "5402"),
		JwtSecret:  getEnvWithDefault(JwtSecret, "secretdb2uy3id28ib3duybc2uy3vfbuyfdkey"),
	}
}

// GetDbDsn constructs the Data Source Name (DSN) for connecting to the database
// using the provided configuration values.
//
// Parameters:
//   - config: a map containing the configuration values needed to construct the DSN.
//
// Returns:
//   - string: the constructed DSN string for the database connection.
func GetDbDsn(config map[string]string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config[DbHost], config[DbUser], config[DbPassword], config[DbName], config[DbPort],
	)
}

// GetJwtSecret retrieves the JWT secret from the configuration map.
//
// Parameters:
//   - config: a map containing the configuration values.
//
// Returns:
//   - string: the JWT secret.
func GetJwtSecret(config map[string]string) string {
	return config[JwtSecret]
}

// GetCorsConfig returns the configuration settings for Cross-Origin Resource Sharing (CORS).
//
// Returns:
//   - *cors.Config: a pointer to the CORS configuration.
func GetCorsConfig() *cors.Config {
	return &cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// DependencyContainer holds the application dependencies for easier access and management.
type DependencyContainer struct {
	DB                    *gorm.DB                  // Database connection
	UserRepository        api.UserRepository        // Repository for user-related operations
	TaskRepository        api.TaskRepository        // Repository for task-related operations
	JwtService            api.JwtService            // Service for handling JWT operations
	PasswordEncoder       api.PasswordEncoder       // Service for password encoding
	AuthenticationService api.AuthenticationService // Service for user authentication
	UserService           api.UserService           // Service for user-related operations
	TaskService           api.TaskService           // Service for task-related operations
	JwtMiddleware         gin.HandlerFunc           // Middleware for JWT authentication
	CorsConfig            *cors.Config              // CORS configuration
}

// NewDependencyContainer initializes and returns a new DependencyContainer.
// It sets up the configuration, database connection, repositories, services, and middleware.
//
// Returns:
//   - *DependencyContainer: a pointer to the newly created DependencyContainer.
//   - error: an error if any initialization fails.
func NewDependencyContainer() (*DependencyContainer, error) {
	config := GetConfig()                // Load configuration from environment variables
	dsn := GetDbDsn(config)              // Get the DSN for database connection
	jwtSecretKey := GetJwtSecret(config) // Get JWT secret key from configuration
	corsConfig := GetCorsConfig()        // Get CORS configuration

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Initialize the database connection
	if err != nil {
		return nil, err // Return an error if the database connection fails
	}

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	taskRepository := repositories.NewTaskRepository(db)

	// Initialize services
	jwtService := services.NewJwtService([]byte(jwtSecretKey))
	passwordEncoder := services.NewDefaultBCryptPasswordEncoder()
	userService := services.NewUserService(userRepository, passwordEncoder)
	taskService := services.NewTaskService(taskRepository, userRepository)
	authenticationService := services.NewAuthenticationService(userService, userRepository, jwtService, passwordEncoder)

	// Initialize middleware
	jwtAuthMiddleware := middlewares.JWTAuthMiddleware(jwtService, userRepository)

	// Return a new DependencyContainer populated with the initialized components
	return &DependencyContainer{
		DB:                    db,
		UserRepository:        userRepository,
		TaskRepository:        taskRepository,
		JwtService:            jwtService,
		PasswordEncoder:       passwordEncoder,
		AuthenticationService: authenticationService,
		UserService:           userService,
		TaskService:           taskService,
		JwtMiddleware:         jwtAuthMiddleware,
		CorsConfig:            corsConfig,
	}, nil
}
