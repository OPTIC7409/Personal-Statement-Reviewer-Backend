package parsing

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExtractJSONToStruct(responseContent string, result interface{}) error {
	startIndex := strings.Index(responseContent, "{")
	if startIndex == -1 {
		return fmt.Errorf("no JSON data found")
	}

	endIndex := strings.LastIndex(responseContent, "}")
	if endIndex == -1 {
		return fmt.Errorf("no closing brace found in response")
	}

	jsonContent := responseContent[startIndex : endIndex+1]

	err := json.Unmarshal([]byte(jsonContent), result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
