package utils

import (
	"encoding/json"
	"strings"
)

// ParseJSONWithCleanup attempts to parse JSON from a string, handling common formatting issues
// like markdown code blocks and extra whitespace
func ParseJSONWithCleanup[T any](data string, target *T) error {
	// First attempt direct parsing
	if err := json.Unmarshal([]byte(data), target); err == nil {
		return nil
	}
	// Try cleaning markdown formatting
	rawResponse := strings.Trim(data, "```json\n")
	rawResponse = strings.Trim(rawResponse, "```\n")
	rawResponse = strings.TrimSpace(rawResponse)

	return json.Unmarshal([]byte(rawResponse), target)
}
