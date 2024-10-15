package repositories

import (
	"github.com/rs/zerolog/log"
	"go_backend/internal/api"
	"go_backend/internal/apperrors"
	"go_backend/internal/models"
	"gorm.io/gorm"
	"strings"
)

// taskRepository is a private struct that implements the TaskRepository interface.
// It provides methods for interacting with tasks in the database.
type taskRepository struct {
	db *gorm.DB // Database connection used to interact with the tasks table.
}

// NewTaskRepository creates a new instance of TaskRepository.
// It initializes a taskRepository with the given database connection.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used for database operations.
//
// Returns:
//   - api.TaskRepository: A new instance of TaskRepository.
func NewTaskRepository(db *gorm.DB) api.TaskRepository {
	log.Debug().Msg("Creating new TaskRepository")
	return &taskRepository{db: db}
}

// FindByTaskID retrieves a task by its TaskID.
//
// Parameters:
//   - taskID: The unique identifier for the task to be retrieved.
//
// Returns:
//   - *models.Task: A pointer to the Task retrieved from the database.
//   - error: An error if the retrieval fails.
func (r *taskRepository) FindByTaskID(taskID int64) (*models.Task, error) {
	if taskID <= 0 {
		return nil, apperrors.NewGenericError("ID is not valid")
	}
	return r.findTask("task_id = ?", taskID)
}

// FindByUserAndName retrieves a task by the user's ID and task name.
//
// Parameters:
//   - user: A pointer to the User struct, which contains the user ID.
//   - name: The name of the task to be retrieved.
//
// Returns:
//   - *models.Task: A pointer to the Task retrieved from the database.
//   - error: An error if the retrieval fails.
func (r *taskRepository) FindByUserAndName(user *models.User, name string) (*models.Task, error) {
	if user == nil || len(strings.TrimSpace(name)) == 0 {
		return nil, apperrors.NewGenericError("Invalid params")
	}
	log.Debug().Int64("userId", user.UserID).Str("taskName", name).Msg("Retrieving task by user and name")
	return r.findTask("user_id = ? AND name = ?", user.UserID, name)
}

// FindByUserAndTaskID retrieves a task by the user's ID and task ID.
//
// Parameters:
//   - user: A pointer to the User struct, which contains the user ID.
//   - taskID: The unique identifier for the task to be retrieved.
//
// Returns:
//   - *models.Task: A pointer to the Task retrieved from the database.
//   - error: An error if the retrieval fails.
func (r *taskRepository) FindByUserAndTaskID(user *models.User, taskID int64) (*models.Task, error) {
	if user == nil || taskID <= 0 {
		return nil, apperrors.NewGenericError("Invalid params")
	}
	log.Debug().Int64("userId", user.UserID).Int64("taskId", taskID).Msg("Retrieving task by user and task ID")
	return r.findTask("user_id = ? AND task_id = ?", user.UserID, taskID)
}

// FindAllByUser retrieves all tasks associated with a specific user.
//
// Parameters:
//   - user: A pointer to the User struct, which contains the user ID.
//
// Returns:
//   - []models.Task: A slice of Task models associated with the user.
//   - error: An error if the retrieval fails.
func (r *taskRepository) FindAllByUser(user *models.User) ([]models.Task, error) {
	if user == nil {
		return nil, apperrors.NewGenericError("User is nil")
	}
	log.Debug().Int64("userId", user.UserID).Msg("Retrieving all tasks for user")
	var tasks []models.Task
	if err := r.db.Where("user_id = ?", user.UserID).Find(&tasks).Error; err != nil {
		log.Error().Err(err).Int64("userId", user.UserID).Msg("Failed to retrieve tasks for user")
		return nil, err
	}
	log.Info().Int64("userId", user.UserID).Int("taskCount", len(tasks)).Msg("Successfully retrieved all tasks for user")
	return tasks, nil
}

// CreateTask creates a new task in the database.
//
// Parameters:
//   - task: A pointer to the Task struct to be created.
//
// Returns:
//   - *models.Task: A pointer to the newly created Task.
//   - error: An error if the creation fails.
func (r *taskRepository) CreateTask(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, apperrors.NewGenericError("Task model is nil")
	}
	log.Debug().Str("taskName", task.Name).Int64("userId", task.UserID).Msg("Creating new task")
	if err := r.db.Create(task).Error; err != nil {
		log.Error().Err(err).Str("taskName", task.Name).Int64("userId", task.UserID).Msg("Failed to create new task")
		return nil, err
	}
	log.Info().Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Successfully created new task")
	return task, nil
}

// UpdateTask updates an existing task in the database.
//
// Parameters:
//   - task: A pointer to the Task struct to be updated.
//
// Returns:
//   - *models.Task: A pointer to the updated Task.
//   - error: An error if the update fails.
func (r *taskRepository) UpdateTask(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, apperrors.NewGenericError("Task model is nil")
	}
	log.Debug().Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Updating task")
	if err := r.db.Save(task).Error; err != nil {
		log.Error().Err(err).Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Failed to update task")
		return nil, err
	}
	log.Info().Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Successfully updated task")
	return task, nil
}

// DeleteTask deletes a task from the database.
//
// Parameters:
//   - task: A pointer to the Task struct to be deleted.
//
// Returns:
//   - error: An error if the deletion fails.
func (r *taskRepository) DeleteTask(task *models.Task) error {
	if task == nil {
		return apperrors.NewGenericError("Task model is nil")
	}
	log.Debug().Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Deleting task")
	if err := r.db.Delete(task).Error; err != nil {
		log.Error().Err(err).Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Failed to delete task")
		return err
	}
	log.Info().Int64("taskId", task.TaskID).Str("taskName", task.Name).Msg("Successfully deleted task")
	return nil
}

// findTask is a private method to encapsulate the task retrieval logic.
// It performs the actual database query and returns the task if found.
//
// Parameters:
//   - query: The SQL query string to execute.
//   - args: The arguments for the SQL query.
//
// Returns:
//   - *models.Task: A pointer to the Task retrieved from the database.
//   - error: An error if the retrieval fails.
func (r *taskRepository) findTask(query string, args ...interface{}) (*models.Task, error) {
	var task models.Task
	result := r.db.Where(query, args...).First(&task)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to retrieve task")
		return nil, result.Error
	}
	log.Info().Int64("taskId", task.TaskID).Msg("Successfully retrieved task")
	return &task, nil
}
