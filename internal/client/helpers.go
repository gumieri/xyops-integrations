package client

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// Helper function for unmarshaling JSON responses.
func unmarshalJSON(body []byte, v any) error {
	if len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, v)
}

// Helper function to check success code in JSON response.
func checkSuccessCode(body []byte) error {
	var resp struct {
		Code int `json:"code"`
	}
	if err := unmarshalJSON(body, &resp); err == nil && resp.Code != 0 {
		var msg struct {
			Message string `json:"message"`
		}
		if err := unmarshalJSON(body, &msg); err != nil {
			log.Printf("unmarshalJSON: failed to extract message from error response: %v", err)
		}
		return NewAPIError(resp.Code, msg.Message, 0)
	}
	return nil
}

// Helper functions for query param building.
func itoa(i int) string {
	return strconv.Itoa(i)
}

func joinComma(vals []string) string {
	return strings.Join(vals, ",")
}
