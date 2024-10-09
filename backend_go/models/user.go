package models

// User represents a user entity in the application.
// It defines the structure of a user and its properties,
// which correspond to the fields in the database table.
//
// Fields:
//   - UserID:       The unique identifier for the user (primary key).
//   - Email:        The user's email address (required and unique).
//   - PasswordHash: The hashed password of the user (required).
//   - Tasks:        A list of tasks associated with the user.
//     The foreign key is UserID, and tasks will be
//     deleted if the user is deleted due to the
//     OnDelete:CASCADE constraint.
type User struct {
	UserID       int64  `gorm:"primaryKey;autoIncrement"`                       // Primary key with auto-increment feature.
	Email        string `gorm:"not null;unique"`                                // User's email; cannot be null and must be unique.
	PasswordHash string `gorm:"not null"`                                       // Hashed password; cannot be null.
	Tasks        []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // Associated tasks, deleted if user is deleted.
}

// NewFullUser creates a new instance of User with the provided details.
// It initializes a User object using the specified parameters.
//
// Parameters:
//   - userId:       The unique identifier for the user.
//   - email:        The email address of the user.
//   - password:     The hashed password of the user.
//   - tasks:        A slice of tasks associated with the user.
//
// Returns:
//   - *User: A pointer to the newly created User instance.
func NewFullUser(userId int64, email, password string, tasks []Task) *User {
	return &User{
		UserID:       userId,
		Email:        email,
		PasswordHash: password,
		Tasks:        tasks,
	}
}

// TableName overrides the default table name for the User struct.
// It specifies the name of the database table to use for the User model.
// In this case, it returns the full name of the table as "backend_diff.users".
//
// Returns:
//   - string: The name of the table as defined in the database.
func (User) TableName() string {
	return "backend_diff.users"
}
