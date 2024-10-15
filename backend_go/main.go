// Package main provides the entry point for the application.
// It initializes logging, sets up dependency injection, configures
// middleware and routes, and starts the HTTP server.
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/internal/config"
	"go_backend/internal/controllers"
	"time"
)

// main is the entry point of the application. It initializes the logger,
// creates the dependency container, sets up controllers, configures the router,
// and starts the server.
func main() {
	// Initialize the logger
	config.InitLogger()
	log.Info().Msg("Starting application")

	// Create the dependency container
	dependencyContainer, dependencyCreationErr := config.NewDependencyContainer()
	if dependencyCreationErr != nil {
		log.Error().Err(dependencyCreationErr).Msg("Failed to create dependency container")
		panic(dependencyContainer)
	}
	log.Info().Msg("Dependency container created successfully")

	// Initialize controllers with their respective services
	authController := controllers.NewAuthController(dependencyContainer.AuthenticationService)
	log.Info().Msg("AuthController initialized")

	userController := controllers.NewUserController(dependencyContainer.UserService)
	log.Info().Msg("UserController initialized")

	taskController := controllers.NewTaskController(dependencyContainer.TaskService)
	log.Info().Msg("TaskController initialized")

	// Set up the Gin router
	router := gin.Default()

	// Middleware to log request details and duration
	router.Use(func(c *gin.Context) {
		start := time.Now()           // Record the start time
		c.Next()                      // Process the request
		duration := time.Since(start) // Calculate duration

		// Log the request details
		log.Debug().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("Request completed")
	})

	// Configure CORS middleware
	router.Use(cors.New(*dependencyContainer.CorsConfig))
	log.Info().Msg("CORS middleware configured")

	// Register routes for authentication, user, and task controllers
	controllers.RegisterAuthRoutes(router, authController)
	log.Info().Msg("Auth routes registered")

	controllers.RegisterUserRoutes(router, userController, dependencyContainer.JwtMiddleware)
	log.Info().Msg("User routes registered")

	controllers.RegisterTaskRoutes(router, taskController, dependencyContainer.JwtMiddleware)
	log.Info().Msg("Task routes registered")

	// Start the server
	log.Info().Msg("Starting server on 0.0.0.0:8080")
	ginAppStartErr := router.Run()
	if ginAppStartErr != nil {
		log.Error().Err(ginAppStartErr).Msg("Failed to start server")
	}
}
