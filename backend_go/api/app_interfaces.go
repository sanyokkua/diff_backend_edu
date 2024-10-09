package api

import (
	"github.com/golang-jwt/jwt/v5"
	"go_backend/dto"
	"go_backend/models"
)

// UserRepository defines the methods for user-related database operations.
type UserRepository interface {
	// CreateUser inserts a new user into the database.
	//
	// Parameters:
	//   - user: a pointer to a models.User struct containing the user information to be inserted.
	//
	// Returns:
	//   - *models.User: a pointer to the newly created user.
	//   - error: an error if the insertion fails.
	CreateUser(user *models.User) (*models.User, error)

	// GetUserByID retrieves a user by their ID.
	//
	// Parameters:
	//   - id: the ID of the user to be retrieved.
	//
	// Returns:
	//   - *models.User: a pointer to the user with the specified ID.
	//   - error: an error if the retrieval fails.
	GetUserByID(id int64) (*models.User, error)

	// UpdateUser updates an existing user's information.
	//
	// Parameters:
	//   - user: a pointer to a models.User struct containing the updated user information.
	//
	// Returns:
	//   - *models.User: a pointer to the updated user.
	//   - error: an error if the update fails.
	UpdateUser(user *models.User) (*models.User, error)

	// DeleteUser removes a user from the database by their ID.
	//
	// Parameters:
	//   - id: the ID of the user to be deleted.
	//
	// Returns:
	//   - error: an error if the deletion fails.
	DeleteUser(id int64) error

	// GetUserByEmail retrieves a user by their email address.
	//
	// Parameters:
	//   - email: the email address of the user to be retrieved.
	//
	// Returns:
	//   - *models.User: a pointer to the user with the specified email address.
	//   - error: an error if the retrieval fails.
	GetUserByEmail(email string) (*models.User, error)
}

// TaskRepository defines the methods for task-related database operations.
type TaskRepository interface {
	// FindByTaskID retrieves a task by its ID.
	//
	// Parameters:
	//   - taskID: the ID of the task to be retrieved.
	//
	// Returns:
	//   - *models.Task: a pointer to the task with the specified ID.
	//   - error: an error if the retrieval fails.
	FindByTaskID(taskID int64) (*models.Task, error)

	// FindByUserAndName retrieves a task by the user and task name.
	//
	// Parameters:
	//   - user: a pointer to a models.User struct representing the user associated with the task.
	//   - name: the name of the task to be retrieved.
	//
	// Returns:
	//   - *models.Task: a pointer to the task matching the user and name.
	//   - error: an error if the retrieval fails.
	FindByUserAndName(user *models.User, name string) (*models.Task, error)

	// FindByUserAndTaskID retrieves a task by the user and task ID.
	//
	// Parameters:
	//   - user: a pointer to a models.User struct representing the user associated with the task.
	//   - taskID: the ID of the task to be retrieved.
	//
	// Returns:
	//   - *models.Task: a pointer to the task matching the user and task ID.
	//   - error: an error if the retrieval fails.
	FindByUserAndTaskID(user *models.User, taskID int64) (*models.Task, error)

	// FindAllByUser retrieves all tasks for a specific user.
	//
	// Parameters:
	//   - user: a pointer to a models.User struct representing the user whose tasks are to be retrieved.
	//
	// Returns:
	//   - []models.Task: a slice containing all tasks associated with the user.
	//   - error: an error if the retrieval fails.
	FindAllByUser(user *models.User) ([]models.Task, error)

	// CreateTask inserts a new task into the database.
	//
	// Parameters:
	//   - task: a pointer to a models.Task struct containing the task information to be inserted.
	//
	// Returns:
	//   - *models.Task: a pointer to the newly created task.
	//   - error: an error if the insertion fails.
	CreateTask(task *models.Task) (*models.Task, error)

	// UpdateTask updates an existing task's information.
	//
	// Parameters:
	//   - task: a pointer to a models.Task struct containing the updated task information.
	//
	// Returns:
	//   - *models.Task: a pointer to the updated task.
	//   - error: an error if the update fails.
	UpdateTask(task *models.Task) (*models.Task, error)

	// DeleteTask removes a task from the database.
	//
	// Parameters:
	//   - task: a pointer to a models.Task struct representing the task to be deleted.
	//
	// Returns:
	//   - error: an error if the deletion fails.
	DeleteTask(task *models.Task) error
}

// JwtService defines the methods for JWT operations.
type JwtService interface {
	// ExtractClaims extracts claims from a JWT token.
	//
	// Parameters:
	//   - token: a string representation of the JWT token from which to extract claims.
	//
	// Returns:
	//   - jwt.Claims: the extracted claims from the token.
	//   - error: an error if the extraction fails.
	ExtractClaims(token string) (jwt.Claims, error)

	// IsTokenExpired checks if the token is expired based on its claims.
	//
	// Parameters:
	//   - claims: the claims extracted from a JWT token.
	//
	// Returns:
	//   - bool: true if the token is expired, false otherwise.
	IsTokenExpired(claims jwt.Claims) bool

	// ValidateToken validates a JWT token against a username.
	//
	// Parameters:
	//   - token: a string representation of the JWT token to validate.
	//   - username: the username to validate against the token.
	//
	// Returns:
	//   - bool: true if the token is valid for the username, false otherwise.
	ValidateToken(token, username string) bool

	// GenerateJwtToken generates a new JWT token for a username.
	//
	// Parameters:
	//   - username: the username for which to generate the token.
	//
	// Returns:
	//   - string: the generated JWT token.
	//   - error: an error if the token generation fails.
	GenerateJwtToken(username string) (string, error)
}

// PasswordEncoder defines the methods for password encoding and matching.
type PasswordEncoder interface {
	// Matches checks if a raw password matches an encoded password.
	//
	// Parameters:
	//   - rawPassword: the raw password provided by the user.
	//   - encodedPassword: the previously encoded password to compare against.
	//
	// Returns:
	//   - bool: true if the passwords match, false otherwise.
	//   - error: an error if the matching process fails.
	Matches(rawPassword, encodedPassword string) (bool, error)

	// Encode encodes a raw password.
	//
	// Parameters:
	//   - rawPassword: the raw password to encode.
	//
	// Returns:
	//   - string: the encoded password.
	//   - error: an error if the encoding fails.
	Encode(rawPassword string) (string, error)
}

// AuthenticationService defines the methods for user authentication.
type AuthenticationService interface {
	// LoginUser authenticates a user and returns their details.
	//
	// Parameters:
	//   - dto: a pointer to a dto.UserLoginDTO containing login credentials.
	//
	// Returns:
	//   - *dto.UserDTO: a pointer to the authenticated user's details.
	//   - error: an error if authentication fails.
	LoginUser(dto *dto.UserLoginDTO) (*dto.UserDTO, error)

	// RegisterUser registers a new user and returns their details.
	//
	// Parameters:
	//   - dto: a pointer to a dto.UserCreationDTO containing the new user's information.
	//
	// Returns:
	//   - *dto.UserDTO: a pointer to the newly registered user's details.
	//   - error: an error if registration fails.
	RegisterUser(dto *dto.UserCreationDTO) (*dto.UserDTO, error)
}

// UserService defines the methods for user-related operations.
type UserService interface {
	// Create creates a new user and returns their details.
	//
	// Parameters:
	//   - userCreationDTO: a pointer to a dto.UserCreationDTO containing the new user's information.
	//
	// Returns:
	//   - *dto.UserDTO: a pointer to the newly created user's details.
	//   - error: an error if creation fails.
	Create(userCreationDTO *dto.UserCreationDTO) (*dto.UserDTO, error)

	// UpdatePassword updates a user's password.
	//
	// Parameters:
	//   - userID: the ID of the user whose password is to be updated.
	//   - userUpdateDTO: a pointer to a dto.UserUpdateDTO containing the new password information.
	//
	// Returns:
	//   - *dto.UserDTO: a pointer to the updated user's details.
	//   - error: an error if the update fails.
	UpdatePassword(userID int64, userUpdateDTO *dto.UserUpdateDTO) (*dto.UserDTO, error)

	// Delete deletes a user by their ID.
	//
	// Parameters:
	//   - userID: the ID of the user to be deleted.
	//   - userDeletionDTO: a pointer to a dto.UserDeletionDTO containing the deletion information.
	//
	// Returns:
	//   - error: an error if the deletion fails.
	Delete(userID int64, userDeletionDTO *dto.UserDeletionDTO) error
}

// TaskService defines the methods for task-related operations.
type TaskService interface {
	// CreateTask creates a new task for a user and returns its details.
	//
	// Parameters:
	//   - userID: the ID of the user for whom the task is being created.
	//   - taskCreationDTO: a pointer to a dto.TaskCreationDTO containing the new task's information.
	//
	// Returns:
	//   - *dto.TaskDTO: a pointer to the newly created task's details.
	//   - error: an error if the creation fails.
	CreateTask(userID int64, taskCreationDTO *dto.TaskCreationDTO) (*dto.TaskDTO, error)

	// UpdateTask updates a task for a user and returns its details.
	//
	// Parameters:
	//   - userID: the ID of the user who owns the task.
	//   - taskID: the ID of the task to be updated.
	//   - taskUpdateDTO: a pointer to a dto.TaskUpdateDTO containing the updated task information.
	//
	// Returns:
	//   - *dto.TaskDTO: a pointer to the updated task's details.
	//   - error: an error if the update fails.
	UpdateTask(userID, taskID int64, taskUpdateDTO *dto.TaskUpdateDTO) (*dto.TaskDTO, error)

	// DeleteTask deletes a task for a user by task ID.
	//
	// Parameters:
	//   - userID: the ID of the user who owns the task.
	//   - taskID: the ID of the task to be deleted.
	//
	// Returns:
	//   - error: an error if the deletion fails.
	DeleteTask(userID, taskID int64) error

	// GetAllTasksForUser retrieves all tasks for a specific user.
	//
	// Parameters:
	//   - userID: the ID of the user whose tasks are to be retrieved.
	//
	// Returns:
	//   - []*dto.TaskDTO: a slice containing pointers to all tasks associated with the user.
	//   - error: an error if the retrieval fails.
	GetAllTasksForUser(userID int64) ([]*dto.TaskDTO, error)

	// GetTaskByUserIDAndTaskID retrieves a task by user ID and task ID.
	//
	// Parameters:
	//   - userID: the ID of the user who owns the task.
	//   - taskID: the ID of the task to be retrieved.
	//
	// Returns:
	//   - *dto.TaskDTO: a pointer to the task associated with the user and task ID.
	//   - error: an error if the retrieval fails.
	GetTaskByUserIDAndTaskID(userID, taskID int64) (*dto.TaskDTO, error)
}
