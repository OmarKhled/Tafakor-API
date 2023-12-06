package utils

import (
	"encoding/json"
	"io"
)

func ParseJSONResponses[Target any](data io.ReadCloser) Target {
	// Final Output
	var res Target

	// Extracting the []byte content
	content, _ := io.ReadAll(data)

	// Checking data json validity
	if valid := json.Valid(content); valid {
		json.Unmarshal(content, &res) // JSON decoding of the body (follows the structure of VersesResponse struct)
	}

	return res
}
