package sdk

import "fmt"

// APIError represents an error returned by the xyops API.
// Per the API reference, the error envelope contains:
// - code: Can be a numeric code (0 for success, non-zero for error) OR a string like "session" for auth errors
// - description: Human-readable error message (not "message")
type APIError struct {
	Code        any    `json:"code"`                  // int or string (e.g., "session")
	Description string `json:"description,omitempty"` // human-readable error message
	HTTPStatus  int
}

func (e *APIError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("API error (code=%v, HTTP=%d): %s", e.Code, e.HTTPStatus, e.Description)
	}
	return fmt.Sprintf("API error (code=%v, HTTP=%d)", e.Code, e.HTTPStatus)
}

// IsAuth returns true if this error is an authentication error.
// The API reports session expiry via code="session" (string) on HTTP 200,
// so we check for that specifically, plus standard HTTP 401/403.
func (e *APIError) IsAuth() bool {
	if e.HTTPStatus == 401 || e.HTTPStatus == 403 {
		return true
	}
	// Check for string "session" code
	if codeStr, ok := e.Code.(string); ok && codeStr == "session" {
		return true
	}
	return false
}

// NewAPIError constructs an APIError from the response.
func NewAPIError(code any, description string, httpStatus int) *APIError {
	return &APIError{
		Code:        code,
		Description: description,
		HTTPStatus:  httpStatus,
	}
}
