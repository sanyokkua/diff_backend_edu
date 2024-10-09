package dto

// ResponseDto is a generic structure for API responses.
//
// It includes the HTTP status code, a status message, and optionally
// the data payload or an error message.
type ResponseDto[T any] struct {
	StatusCode    int    `json:"statusCode"`      // HTTP status code for the response
	StatusMessage string `json:"statusMessage"`   // Human-readable message describing the status
	Data          T      `json:"data,omitempty"`  // Actual data payload (omitted if nil)
	Error         string `json:"error,omitempty"` // Error message (omitted if empty)
}

// NewFullResponseDto creates a new ResponseDto instance with the provided parameters.
//
// Parameters:
//   - statusCode: the HTTP status code to return.
//   - statusMessage: a message describing the status.
//   - data: the data to include in the response. Can be of any type.
//   - err: an optional error message.
//
// Returns:
//   - *ResponseDto[any]: a pointer to the newly created ResponseDto instance.
func NewFullResponseDto(statusCode int, statusMessage string, data any, err string) *ResponseDto[any] {
	return &ResponseDto[any]{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Data:          data,
		Error:         err,
	}
}
