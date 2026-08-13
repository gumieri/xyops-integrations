package sdk

import "encoding/json"

// Helper function for unmarshaling JSON responses.
// Returns nil for empty bodies (e.g., 204 No Content), which callers interpret
// as success with no payload to unmarshal. v must be a non-nil pointer.
func unmarshalJSON(body []byte, v any) error {
	if len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, v)
}
