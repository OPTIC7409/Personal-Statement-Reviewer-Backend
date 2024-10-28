package aidetection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"psr/types/aidetection"
	"strings"
)

func DetectAIContent(personalStatement string) (aidetection.AIDetectionResult, error) {
	url := "https://api.gptzero.me/v2/predict/text"
	sections := strings.Split(personalStatement, "\n\n")

	var flaggedSections []aidetection.FlaggedSection
	var totalProbability float64

	for _, section := range sections {
		payload := map[string]interface{}{
			"document":     section,
			"multilingual": false,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return aidetection.AIDetectionResult{}, fmt.Errorf("Error marshaling JSON: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return aidetection.AIDetectionResult{}, fmt.Errorf("HTTP request error: %v", err)
		}
		defer resp.Body.Close()

		var result struct {
			Documents []struct {
				AverageGeneratedProb float64 `json:"average_generated_prob"`
			}
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return aidetection.AIDetectionResult{}, fmt.Errorf("Error parsing JSON response: %v", err)
		}

		if len(result.Documents) > 0 {
			probability := result.Documents[0].AverageGeneratedProb
			totalProbability += probability

			if probability > 50 {
				flaggedSections = append(flaggedSections, aidetection.FlaggedSection{
					Text:        section,
					Reason:      "High AI probability",
					Probability: probability,
				})
			}
		}
	}

	overallProbability := totalProbability / float64(len(sections))

	detectionResult := aidetection.AIDetectionResult{
		OverallAIProbability: overallProbability,
		FlaggedSections:      flaggedSections,
	}

	return detectionResult, nil
}
