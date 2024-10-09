package dto

// UserCreationDTO is used for creating a new user.
// It includes the user's email, password, and a confirmation password
// for validation.
type UserCreationDTO struct {
	Email                string `json:"email" binding:"required,email"`          // Required: User's email, must be a valid email format
	Password             string `json:"password" binding:"required,min=4"`       // Required: User's password, minimum 4 characters
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"` // Required: Confirmation of the user's password
}

// UserDeletionDTO is used for deleting a user.
// It includes the user's email and the current password for
// verification before deletion.
type UserDeletionDTO struct {
	Email           string `json:"email" binding:"required,email"`     // Required: User's email, must be a valid email format
	CurrentPassword string `json:"currentPassword" binding:"required"` // Required: User's current password for authentication
}

// UserDTO represents a user in the system.
// It includes the user ID, email, and a JWT token for authentication.
type UserDTO struct {
	UserID   int64  `json:"userId"`   // Unique identifier for the user
	Email    string `json:"email"`    // User's email address
	JwtToken string `json:"jwtToken"` // JWT token for user authentication
}

// NewFullUserDto creates a new UserDTO instance.
// This constructor is useful for creating a UserDTO with user ID, email, and JWT token.
func NewFullUserDto(userID int64, email, jwtToken string) *UserDTO {
	return &UserDTO{
		UserID:   userID,
		Email:    email,
		JwtToken: jwtToken,
	}
}

// UserLoginDTO is used for user login requests.
// It includes the user's email and password for authentication.
type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"` // Required: User's email, must be a valid email format
	Password string `json:"password" binding:"required"`    // Required: User's password for authentication
}

// UserUpdateDTO is used for updating a user's password.
// It includes the current password, new password, and a confirmation
// of the new password.
type UserUpdateDTO struct {
	CurrentPassword         string `json:"currentPassword" binding:"required"`         // Required: User's current password
	NewPassword             string `json:"newPassword" binding:"required,min=8"`       // Required: New password, minimum 8 characters
	NewPasswordConfirmation string `json:"newPasswordConfirmation" binding:"required"` // Required: Confirmation of the new password
}
