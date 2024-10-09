package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/dto"
	"go_backend/models"
	"go_backend/utils"
	"net/http"
)

// TaskController handles task-related operations, such as creating, retrieving, updating, and deleting tasks.
type TaskController struct {
	taskService api.TaskService // The service used for task operations
}

// NewTaskController initializes a new TaskController instance.
//
// Parameters:
//   - taskService: an instance of TaskService for handling task operations.
//
// Returns:
//   - *TaskController: a pointer to the newly created TaskController instance.
func NewTaskController(taskService api.TaskService) *TaskController {
	log.Debug().Msg("Initializing TaskController")
	return &TaskController{
		taskService: taskService,
	}
}

// RegisterTaskRoutes registers task-related routes with the provided router.
//
// Parameters:
//   - router: the Gin router to which the task routes will be registered.
//   - taskController: the TaskController instance that handles task logic.
//   - middleware: a middleware function to apply to the task routes.
func RegisterTaskRoutes(router *gin.Engine, taskController *TaskController, middleware gin.HandlerFunc) {
	v1 := router.Group("api/v1/users/:userId/tasks") // Grouping routes under user tasks
	v1.Use(middleware)                               // Apply middleware for authentication
	{
		v1.POST("/", taskController.createTask)          // Route to create a new task
		v1.GET("/:taskId", taskController.getTaskByID)   // Route to retrieve a task by ID
		v1.PUT("/:taskId", taskController.updateTask)    // Route to update an existing task
		v1.DELETE("/:taskId", taskController.deleteTask) // Route to delete a task
		v1.GET("/", taskController.getAllTasksForUser)   // Route to get all tasks for a user
	}
}

// extractAndValidateUserTask handles repeated logic for extracting and validating user and task IDs.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// Returns:
//   - userID: the ID of the user extracted from the request context.
//   - userFromContext: the user object extracted from the context.
//   - taskID: the ID of the task extracted from the request.
//   - error: any error encountered during extraction and validation.
func (r *TaskController) extractAndValidateUserTask(ctx *gin.Context) (int64, *models.User, int64, error) {
	userFromContext, err := getUserFromContext(ctx)
	if err != nil {
		return 0, nil, 0, err
	}

	userID, err := getIntParamByName(ctx, "userId")
	if err != nil {
		return 0, nil, 0, err
	}

	if err := utils.ValidateAuthenticatedUserID(userFromContext.UserID, userID); err != nil {
		return 0, nil, 0, err
	}

	taskID, err := getIntParamByName(ctx, "taskId")
	if err != nil {
		return 0, nil, 0, err
	}

	return userID, userFromContext, taskID, nil
}

// createTask handles the creation of a new task.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function attempts to bind the incoming JSON request body to a TaskCreationDTO,
// validates the user, and calls the task service to create a new task. If successful,
// it returns the newly created task; otherwise, it responds with an error.
func (r *TaskController) createTask(ctx *gin.Context) {
	log.Debug().Msg("Handling createTask request")

	userFromContext, err := getUserFromContext(ctx)
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	userID, err := getIntParamByName(ctx, "userId")
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := utils.ValidateAuthenticatedUserID(userFromContext.UserID, userID); err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var taskCreationDTO dto.TaskCreationDTO
	if err := ctx.BindJSON(&taskCreationDTO); err != nil {
		log.Error().Err(err).Msg("Failed to parse TaskCreationDTO")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	newTask, err := r.taskService.CreateTask(userID, &taskCreationDTO)
	if err != nil {
		log.Error().Err(err).Int64("userId", userID).Msg("Task creation failed")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("userId", userID).Msg("Task created successfully")
	utils.WriteSuccessResponse(ctx, http.StatusCreated, newTask) // Return success response with created task
}

// getTaskByID handles retrieval of a task by its ID.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user and task IDs, then retrieves the task from the task service.
// If successful, it returns the task; otherwise, it responds with an error.
func (r *TaskController) getTaskByID(ctx *gin.Context) {
	log.Debug().Msg("Handling getTaskByID request")

	userID, _, taskID, err := r.extractAndValidateUserTask(ctx)
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	task, err := r.taskService.GetTaskByUserIDAndTaskID(userID, taskID)
	if err != nil {
		log.Error().Err(err).Int64("taskId", taskID).Msg("Task retrieval failed")
		handleErrorResponse(ctx, err, http.StatusNotFound)
		return
	}

	log.Info().Int64("taskId", taskID).Msg("Task retrieved successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, task) // Return success response with task details
}

// updateTask handles updating an existing task.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user and task IDs, binds the request body to a TaskUpdateDTO,
// and calls the task service to update the task. If successful, it returns the updated task;
// otherwise, it responds with an error.
func (r *TaskController) updateTask(ctx *gin.Context) {
	log.Debug().Msg("Handling updateTask request")

	userID, _, taskID, err := r.extractAndValidateUserTask(ctx)
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	var taskUpdateDTO dto.TaskUpdateDTO
	if err := ctx.BindJSON(&taskUpdateDTO); err != nil {
		log.Error().Err(err).Msg("Failed to parse TaskUpdateDTO")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	updatedTask, err := r.taskService.UpdateTask(userID, taskID, &taskUpdateDTO)
	if err != nil {
		log.Error().Err(err).Int64("taskId", taskID).Msg("Task update failed")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("taskId", taskID).Msg("Task updated successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, updatedTask) // Return success response with updated task
}

// deleteTask handles the deletion of a task.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user and task IDs, then calls the task service to delete the task.
// If successful, it responds with a no content status; otherwise, it responds with an error.
func (r *TaskController) deleteTask(ctx *gin.Context) {
	log.Debug().Msg("Handling deleteTask request")

	userID, _, taskID, err := r.extractAndValidateUserTask(ctx)
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	if err := r.taskService.DeleteTask(userID, taskID); err != nil {
		log.Error().Err(err).Int64("taskId", taskID).Msg("Task deletion failed")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("taskId", taskID).Msg("Task deleted successfully")
	utils.WriteSuccessResponse(ctx, http.StatusNoContent, nil) // Respond with no content on successful deletion
}

// getAllTasksForUser retrieves all tasks associated with a specific user.
//
// Parameters:
//   - ctx: the Gin context, which carries the request and response data.
//
// This function validates the user ID, retrieves all tasks for the user from the task service,
// and responds with the list of tasks or an error if encountered.
func (r *TaskController) getAllTasksForUser(ctx *gin.Context) {
	log.Debug().Msg("Handling getAllTasksForUser request")

	userFromContext, err := getUserFromContext(ctx)
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	userID, err := getIntParamByName(ctx, "userId")
	if err != nil {
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := utils.ValidateAuthenticatedUserID(userFromContext.UserID, userID); err != nil {
		handleErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	tasks, err := r.taskService.GetAllTasksForUser(userID)
	if err != nil {
		log.Error().Err(err).Int64("userId", userID).Msg("Failed to retrieve tasks")
		handleErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Info().Int64("userId", userID).Msg("All tasks retrieved successfully")
	utils.WriteSuccessResponse(ctx, http.StatusOK, tasks) // Return success response with list of tasks
}
