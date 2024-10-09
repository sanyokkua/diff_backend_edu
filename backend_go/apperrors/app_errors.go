package apperrors

import "errors"

type AppError interface {
	StatusCode() int
	Error() string
}

type EmailAlreadyExistsError struct {
	ErrorMsg string
}

func (r EmailAlreadyExistsError) Error() string {
	if r.ErrorMsg == "" {
		return "Email is already in use"
	}
	return r.ErrorMsg
}

type InvalidEmailFormatError struct {
	ErrorMsg string
}

func (r InvalidEmailFormatError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid email format"
	}
	return r.ErrorMsg
}

type InvalidJwtTokenError struct {
	ErrorMsg string
}

func (r InvalidJwtTokenError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid JWT token"
	}
	return r.ErrorMsg
}

type InvalidPasswordError struct {
	ErrorMsg string
}

func (r InvalidPasswordError) Error() string {
	if r.ErrorMsg == "" {
		return "Invalid credentials"
	}
	return r.ErrorMsg
}

type TaskAlreadyExistsError struct {
	ErrorMsg string
}

func (r TaskAlreadyExistsError) Error() string {
	if r.ErrorMsg == "" {
		return "Task with the this name already exists for the user"
	}
	return r.ErrorMsg
}

type TaskNotFoundError struct {
	ErrorMsg string
}

func (r TaskNotFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "Task not found"
	}
	return r.ErrorMsg
}

type AccessDeniedError struct {
	ErrorMsg string
}

func (r AccessDeniedError) Error() string {
	if r.ErrorMsg == "" {
		return "Access Denied"
	}
	return r.ErrorMsg
}

type AuthenticationCredentialsNotFoundError struct {
	ErrorMsg string
}

func (r AuthenticationCredentialsNotFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "Authentication Credentials Not Found"
	}
	return r.ErrorMsg
}

type InsufficientAuthenticationError struct {
	ErrorMsg string
}

func (r InsufficientAuthenticationError) Error() string {
	if r.ErrorMsg == "" {
		return "Insufficient Authentication"
	}
	return r.ErrorMsg
}

type NoHandlerFoundError struct {
	ErrorMsg string
}

func (r NoHandlerFoundError) Error() string {
	if r.ErrorMsg == "" {
		return "No Handler Found"
	}
	return r.ErrorMsg
}

type IllegalArgumentError struct {
	ErrorMsg string
}

func (r IllegalArgumentError) Error() string {
	if r.ErrorMsg == "" {
		return "Illegal Argument Error"
	}
	return r.ErrorMsg
}

func NewGenericError(msg string) error {
	return errors.New(msg)
}
func NewEmailAlreadyExistsError(msg string) EmailAlreadyExistsError {
	return EmailAlreadyExistsError{ErrorMsg: msg}
}
func NewInvalidEmailFormatError(msg string) InvalidEmailFormatError {
	return InvalidEmailFormatError{ErrorMsg: msg}
}
func NewInvalidJwtTokenError(msg string) InvalidJwtTokenError {
	return InvalidJwtTokenError{ErrorMsg: msg}
}
func NewInvalidPasswordError(msg string) InvalidPasswordError {
	return InvalidPasswordError{ErrorMsg: msg}
}
func NewTaskAlreadyExistsError(msg string) TaskAlreadyExistsError {
	return TaskAlreadyExistsError{ErrorMsg: msg}
}
func NewTaskNotFoundError(msg string) TaskNotFoundError {
	return TaskNotFoundError{ErrorMsg: msg}
}
func NewAccessDeniedError(msg string) AccessDeniedError {
	return AccessDeniedError{ErrorMsg: msg}
}
func NewIllegalArgumentError(msg string) IllegalArgumentError {
	return IllegalArgumentError{ErrorMsg: msg}
}

var EmailAlreadyExistsErrorDefault = NewEmailAlreadyExistsError("")
var InvalidEmailFormatErrorDefault = NewInvalidEmailFormatError("")
