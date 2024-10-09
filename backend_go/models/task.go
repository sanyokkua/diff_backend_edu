package models

// Task represents a task entity in the application.
// It defines the structure of a task and its properties,
// which correspond to the fields in the database table.
//
// Fields:
//   - TaskID:      The unique identifier for the task (primary key).
//   - Name:        The name of the task (required).
//   - Description: A detailed description of the task (required).
//   - UserID:      The ID of the user associated with the task (required).
type Task struct {
	TaskID      int64  `gorm:"primaryKey;autoIncrement"` // Primary key with auto-increment feature.
	Name        string `gorm:"not null"`                 // Task name; cannot be null.
	Description string `gorm:"not null"`                 // Task description; cannot be null.
	UserID      int64  `gorm:"not null"`                 // User ID associated with the task; cannot be null.
}

// TableName overrides the default table name for the Task struct.
// It specifies the name of the database table to use for the Task model.
// In this case, it returns the full name of the table as "backend_diff.tasks".
//
// Returns:
//   - string: The name of the table as defined in the database.
func (Task) TableName() string {
	return "backend_diff.tasks"
}
