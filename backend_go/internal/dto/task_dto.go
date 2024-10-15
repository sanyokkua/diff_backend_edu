package dto

// TaskDTO represents a task in the system with all its properties.
//
// It includes the task ID, name, description, and the ID of the user
// who owns the task.
type TaskDTO struct {
	TaskID      int64  `json:"taskId"`      // Unique identifier for the task
	Name        string `json:"name"`        // Name of the task
	Description string `json:"description"` // Detailed description of the task
	UserID      int64  `json:"userId"`      // ID of the user who owns the task
}

// TaskCreationDTO is used for creating a new task.
//
// It contains the necessary fields to create a new task without exposing
// the TaskID or UserID, which are managed by the system.
type TaskCreationDTO struct {
	Name        string `json:"name" binding:"required"`        // Name of the task, required field
	Description string `json:"description" binding:"required"` // Description of the task, required field
}

// TaskUpdateDTO is used for updating an existing task.
//
// It includes the fields that can be modified for a task, allowing for partial updates.
type TaskUpdateDTO struct {
	Name        string `json:"name" binding:"omitempty"`        // Optional: New name of the task
	Description string `json:"description" binding:"omitempty"` // Optional: New description of the task
}
