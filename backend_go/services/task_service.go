package services

import (
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/apperrors"
	"go_backend/dto"
	"go_backend/models"
	"go_backend/utils"
)

// taskService implements TaskService for managing tasks.
type taskService struct {
	taskRepo api.TaskRepository // Repository for task data
	userRepo api.UserRepository // Repository for user data
}

// NewTaskService creates a new TaskService with the provided repositories.
//
// Parameters:
//   - taskRepo: An api.TaskRepository instance for task operations.
//   - userRepo: An api.UserRepository instance for user operations.
//
// Returns:
//   - api.TaskService: A new instance of TaskService.
func NewTaskService(taskRepo api.TaskRepository, userRepo api.UserRepository) api.TaskService {
	log.Debug().Msg("Creating new TaskService instance")
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// CreateTask creates a new task for the specified user.
//
// Parameters:
//   - userId: The ID of the user creating the task.
//   - taskCreationDTO: A pointer to dto.TaskCreationDTO containing task details.
//
// Returns:
//   - *dto.TaskDTO: The created task as a TaskDTO.
//   - error: An error if task creation fails.
func (s *taskService) CreateTask(userId int64, taskCreationDTO *dto.TaskCreationDTO) (*dto.TaskDTO, error) {
	log.Debug().Int64("userId", userId).Str("taskName", taskCreationDTO.Name).Msg("Entering CreateTask")

	// Validate task data
	if err := utils.ValidateTaskCreation(taskCreationDTO); err != nil {
		log.Warn().Err(err).Int64("userId", userId).Str("taskName", taskCreationDTO.Name).Msg("Task validation failed")
		return nil, err
	}

	// Retrieve and validate user
	user, err := s.getUserByID(userId)
	if err != nil {
		return nil, err
	}

	// Ensure task name is unique for the user
	if err := s.ensureUniqueTaskName(user, taskCreationDTO.Name); err != nil {
		return nil, err
	}

	// Create new task
	newTask := &models.Task{
		Name:        taskCreationDTO.Name,
		Description: taskCreationDTO.Description,
		UserID:      user.UserID,
	}

	// Save the task
	return s.saveTask(newTask)
}

// UpdateTask updates an existing task for a user.
//
// Parameters:
//   - userId: The ID of the user updating the task.
//   - taskId: The ID of the task to be updated.
//   - taskUpdateDTO: A pointer to dto.TaskUpdateDTO containing updated task details.
//
// Returns:
//   - *dto.TaskDTO: The updated task as a TaskDTO.
//   - error: An error if task update fails.
func (s *taskService) UpdateTask(userId, taskId int64, taskUpdateDTO *dto.TaskUpdateDTO) (*dto.TaskDTO, error) {
	log.Debug().Int64("userId", userId).Int64("taskId", taskId).Msg("Entering UpdateTask")

	// Validate task data
	if err := utils.ValidateTaskUpdateDTO(taskUpdateDTO); err != nil {
		log.Warn().Err(err).Int64("userId", userId).Int64("taskId", taskId).Msg("Task validation failed")
		return nil, err
	}

	// Fetch and verify task ownership
	task, err := s.getTaskForUser(userId, taskId)
	if err != nil {
		return nil, err
	}

	// Update task details
	task.Name = taskUpdateDTO.Name
	task.Description = taskUpdateDTO.Description

	// Save updated task
	return s.saveTask(task)
}

// DeleteTask removes a task if it belongs to the specified user.
//
// Parameters:
//   - userId: The ID of the user requesting deletion.
//   - taskId: The ID of the task to be deleted.
//
// Returns:
//   - error: An error if task deletion fails.
func (s *taskService) DeleteTask(userId, taskId int64) error {
	log.Debug().Int64("userId", userId).Int64("taskId", taskId).Msg("Entering DeleteTask")

	// Fetch and verify task ownership
	task, err := s.getTaskForUser(userId, taskId)
	if err != nil {
		return err
	}

	// Delete the task
	if err := s.taskRepo.DeleteTask(task); err != nil {
		log.Error().Err(err).Int64("userId", userId).Int64("taskId", taskId).Msg("Error deleting task")
		return apperrors.NewGenericError(err.Error())
	}

	log.Info().Int64("userId", userId).Int64("taskId", taskId).Msg("Task deleted successfully")
	return nil
}

// GetAllTasksForUser retrieves all tasks for a specific user.
//
// Parameters:
//   - userId: The ID of the user whose tasks are to be retrieved.
//
// Returns:
//   - []*dto.TaskDTO: A slice of TaskDTOs representing the user's tasks.
//   - error: An error if the retrieval fails.
func (s *taskService) GetAllTasksForUser(userId int64) ([]*dto.TaskDTO, error) {
	log.Debug().Int64("userId", userId).Msg("Entering GetAllTasksForUser")

	// Retrieve user
	user, err := s.getUserByID(userId)
	if err != nil {
		return nil, err
	}

	// Retrieve tasks for the user
	tasks, err := s.taskRepo.FindAllByUser(user)
	if err != nil {
		log.Error().Err(err).Int64("userId", userId).Msg("Error retrieving tasks for user")
		return nil, err
	}

	// Convert to DTOs
	return s.convertToTaskDTOs(tasks), nil
}

// GetTaskByUserIDAndTaskID retrieves a specific task by ID for a user.
//
// Parameters:
//   - userId: The ID of the user requesting the task.
//   - taskId: The ID of the task to retrieve.
//
// Returns:
//   - *dto.TaskDTO: The requested task as a TaskDTO.
//   - error: An error if the task retrieval fails.
func (s *taskService) GetTaskByUserIDAndTaskID(userId, taskId int64) (*dto.TaskDTO, error) {
	log.Debug().Int64("userId", userId).Int64("taskId", taskId).Msg("Entering GetTaskByUserIDAndTaskID")

	// Retrieve user
	user, err := s.getUserByID(userId)
	if err != nil {
		return nil, err
	}

	// Retrieve task by user and task ID
	task, err := s.taskRepo.FindByUserAndTaskID(user, taskId)
	if err != nil {
		log.Error().Err(err).Int64("userId", userId).Int64("taskId", taskId).Msg("Task not found for user")
		return nil, apperrors.NewTaskNotFoundError(err.Error())
	}

	// Convert to DTO and return
	return s.convertToTaskDTO(task), nil
}

// Helper methods

// getUserByID retrieves a user by ID and logs any errors.
//
// Parameters:
//   - userId: The ID of the user to retrieve.
//
// Returns:
//   - *models.User: The retrieved user.
//   - error: An error if user retrieval fails.
func (s *taskService) getUserByID(userId int64) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userId)
	if err != nil {
		log.Error().Err(err).Int64("userId", userId).Msg("Error finding user")
		return nil, err
	}
	return user, nil
}

// getTaskForUser ensures the task exists and belongs to the user.
//
// Parameters:
//   - userId: The ID of the user requesting the task.
//   - taskId: The ID of the task to verify.
//
// Returns:
//   - *models.Task: The verified task.
//   - error: An error if the task verification fails.
func (s *taskService) getTaskForUser(userId, taskId int64) (*models.Task, error) {
	task, err := s.taskRepo.FindByTaskID(taskId)
	if err != nil {
		log.Error().Err(err).Int64("taskId", taskId).Msg("Task not found")
		return nil, apperrors.NewTaskNotFoundError(err.Error())
	}
	if task.UserID != userId {
		log.Warn().Int64("userId", userId).Int64("taskId", taskId).Msg("Task does not belong to the user")
		return nil, apperrors.NewAccessDeniedError("Task does not belong to the user")
	}
	return task, nil
}

// ensureUniqueTaskName checks if a task with the same name exists for the user.
//
// Parameters:
//   - user: The user to check for existing tasks.
//   - taskName: The name of the task to check.
//
// Returns:
//   - error: An error if a task with the same name already exists.
func (s *taskService) ensureUniqueTaskName(user *models.User, taskName string) error {
	if err := utils.CheckTaskExistsForUser(s.taskRepo, user, taskName); err != nil {
		log.Warn().Err(err).Int64("userId", user.UserID).Str("taskName", taskName).Msg("Task with the same name already exists")
		return err
	}
	return nil
}

// saveTask saves the task and converts it to a DTO.
//
// Parameters:
//   - task: The task to save.
//
// Returns:
//   - *dto.TaskDTO: The saved task as a TaskDTO.
//   - error: An error if saving the task fails.
func (s *taskService) saveTask(task *models.Task) (*dto.TaskDTO, error) {
	savedTask, err := s.taskRepo.CreateTask(task)
	if err != nil {
		log.Error().Err(err).Int64("taskId", task.TaskID).Msg("Error saving task")
		return nil, err
	}

	log.Info().Int64("taskId", savedTask.TaskID).Str("taskName", savedTask.Name).Msg("Task saved successfully")
	return s.convertToTaskDTO(savedTask), nil
}

// convertToTaskDTO converts a task to a TaskDTO.
//
// Parameters:
//   - task: The task to convert.
//
// Returns:
//   - *dto.TaskDTO: The converted TaskDTO.
func (s *taskService) convertToTaskDTO(task *models.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		TaskID:      task.TaskID,
		Name:        task.Name,
		Description: task.Description,
		UserID:      task.UserID,
	}
}

// convertToTaskDTOs converts a list of tasks to a list of TaskDTOs.
//
// Parameters:
//   - tasks: A slice of tasks to convert.
//
// Returns:
//   - []*dto.TaskDTO: A slice of converted TaskDTOs.
func (s *taskService) convertToTaskDTOs(tasks []models.Task) []*dto.TaskDTO {
	taskDTOs := make([]*dto.TaskDTO, len(tasks))
	for i, task := range tasks {
		taskDTOs[i] = s.convertToTaskDTO(&task)
	}
	return taskDTOs
}
