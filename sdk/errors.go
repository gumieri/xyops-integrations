package sdk

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
// The API reports session expiry inconsistently: sometimes via HTTP status
// (401/403) and sometimes via a "session" message on a 200 response, so both
// are checked.
func (e *APIError) IsAuth() bool {
	return e.HTTPStatus == 401 || e.HTTPStatus == 403 || e.Message == "session"
}

// NewAPIError constructs an APIError from the response.
func NewAPIError(code int, message string, httpStatus int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}
