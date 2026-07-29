package client

import "fmt"

// APIError represents an error returned by the xyops API.
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message,omitempty"`
	HTTPStatus int
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error (code=%d, HTTP=%d): %s", e.Code, e.HTTPStatus, e.Message)
	}
	return fmt.Sprintf("API error (code=%d, HTTP=%d)", e.Code, e.HTTPStatus)
}

// IsAuth returns true if this error is an authentication error.
func (e *APIError) IsAuth() bool {
	return (e.Code != 0 || e.Message == "session") && (e.Message == "session" || e.HTTPStatus == 401 || e.HTTPStatus == 403)
}

// NewAPIError constructs an APIError from the response.
func NewAPIError(code int, message string, httpStatus int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}
