// Package apperrors provides application-specific error types.
package apperrors

import "errors"

// EmailAlreadyExistsError represents an error when the email is already in use.
type EmailAlreadyExistsError struct {
	ErrorMsg string
}

// Error returns the error message for EmailAlreadyExistsError.
func (r EmailAlreadyExistsError) Error() string {
	if r.ErrorMsg == "" {
		return "Email is already in use"
	}
	return r.ErrorMsg
}

// InvalidEmailFormatError represents an error when the email format is invalid.
type InvalidEmailFormatError struct {
	ErrorMsg string
}

// Error returns the error message for InvalidEmailFormatError.
func (r InvalidEmailFormatError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid email format"
	}
	return r.ErrorMsg
}

// InvalidJwtTokenError represents an error when the JWT token is invalid.
type InvalidJwtTokenError struct {
	ErrorMsg string
}

// Error returns the error message for InvalidJwtTokenError.
func (r InvalidJwtTokenError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid JWT token"
	}
	return r.ErrorMsg
}

// InvalidPasswordError represents an error when the password is invalid.
type InvalidPasswordError struct {
	ErrorMsg string
}

// Error returns the error message for InvalidPasswordError.
func (r InvalidPasswordError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid credentials"
	}
	return r.ErrorMsg
}

// TaskAlreadyExistsError represents an error when a task with the same name already exists for the user.
type TaskAlreadyExistsError struct {
	ErrorMsg string
}

// Error returns the error message for TaskAlreadyExistsError.
func (r TaskAlreadyExistsError) Error() string {
	if r.ErrorMsg == "" {
		return "Task with this name already exists for the user"
	}
	return r.ErrorMsg
}

// TaskNotFoundError represents an error when a task is not found.
type TaskNotFoundError struct {
	ErrorMsg string
}

// Error returns the error message for TaskNotFoundError.
func (r TaskNotFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "Task not found"
	}
	return r.ErrorMsg
}

// AccessDeniedError represents an error when access is denied.
type AccessDeniedError struct {
	ErrorMsg string
}

// Error returns the error message for AccessDeniedError.
func (r AccessDeniedError) Error() string {
	if r.ErrorMsg == "" {
		return "Access Denied"
	}
	return r.ErrorMsg
}

// AuthenticationCredentialsNotFoundError represents an error when authentication credentials are not found.
type AuthenticationCredentialsNotFoundError struct {
	ErrorMsg string
}

// Error returns the error message for AuthenticationCredentialsNotFoundError.
func (r AuthenticationCredentialsNotFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "Authentication Credentials Not Found"
	}
	return r.ErrorMsg
}

// InsufficientAuthenticationError represents an error when authentication is insufficient.
type InsufficientAuthenticationError struct {
	ErrorMsg string
}

// Error returns the error message for InsufficientAuthenticationError.
func (r InsufficientAuthenticationError) Error() string {
	if r.ErrorMsg == "" {
		return "Insufficient Authentication"
	}
	return r.ErrorMsg
}

// NoHandlerFoundError represents an error when no handler is found.
type NoHandlerFoundError struct {
	ErrorMsg string
}

// Error returns the error message for NoHandlerFoundError.
func (r NoHandlerFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "No Handler Found"
	}
	return r.ErrorMsg
}

// IllegalArgumentError represents an error when an illegal argument is provided.
type IllegalArgumentError struct {
	ErrorMsg string
}

// Error returns the error message for IllegalArgumentError.
func (r IllegalArgumentError) Error() string {
	if r.ErrorMsg == "" {
		return "Illegal Argument Error"
	}
	return r.ErrorMsg
}

// NewGenericError creates a new generic error with the given message.
func NewGenericError(msg string) error {
	return errors.New(msg)
}

// NewEmailAlreadyExistsError creates a new EmailAlreadyExistsError with the given message.
func NewEmailAlreadyExistsError(msg string) EmailAlreadyExistsError {
	return EmailAlreadyExistsError{ErrorMsg: msg}
}

// NewInvalidEmailFormatError creates a new InvalidEmailFormatError with the given message.
func NewInvalidEmailFormatError(msg string) InvalidEmailFormatError {
	return InvalidEmailFormatError{ErrorMsg: msg}
}

// NewInvalidJwtTokenError creates a new InvalidJwtTokenError with the given message.
func NewInvalidJwtTokenError(msg string) InvalidJwtTokenError {
	return InvalidJwtTokenError{ErrorMsg: msg}
}

// NewInvalidPasswordError creates a new InvalidPasswordError with the given message.
func NewInvalidPasswordError(msg string) InvalidPasswordError {
	return InvalidPasswordError{ErrorMsg: msg}
}

// NewTaskAlreadyExistsError creates a new TaskAlreadyExistsError with the given message.
func NewTaskAlreadyExistsError(msg string) TaskAlreadyExistsError {
	return TaskAlreadyExistsError{ErrorMsg: msg}
}

// NewTaskNotFoundError creates a new TaskNotFoundError with the given message.
func NewTaskNotFoundError(msg string) TaskNotFoundError {
	return TaskNotFoundError{ErrorMsg: msg}
}

// NewAccessDeniedError creates a new AccessDeniedError with the given message.
func NewAccessDeniedError(msg string) AccessDeniedError {
	return AccessDeniedError{ErrorMsg: msg}
}

// NewIllegalArgumentError creates a new IllegalArgumentError with the given message.
func NewIllegalArgumentError(msg string) IllegalArgumentError {
	return IllegalArgumentError{ErrorMsg: msg}
}

// NewAuthenticationCredentialsNotFoundError creates a new AuthenticationCredentialsNotFoundError with the given message.
func NewAuthenticationCredentialsNotFoundError(msg string) AuthenticationCredentialsNotFoundError {
	return AuthenticationCredentialsNotFoundError{ErrorMsg: msg}
}

// NewInsufficientAuthenticationError creates a new InsufficientAuthenticationError with the given message.
func NewInsufficientAuthenticationError(msg string) InsufficientAuthenticationError {
	return InsufficientAuthenticationError{ErrorMsg: msg}
}

var (
	// EmailAlreadyExistsErrorDefault is the default EmailAlreadyExistsError.
	EmailAlreadyExistsErrorDefault = NewEmailAlreadyExistsError("")

	// InvalidEmailFormatErrorDefault is the default InvalidEmailFormatError.
	InvalidEmailFormatErrorDefault = NewInvalidEmailFormatError("")
)
