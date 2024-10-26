package winston

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"psr/types/aidetection"
	"psr/types/plagiarism"
)

const (
	apiKey  = "wm8DkvVUHLkEQlR1V6J2GGVy5guHbzhjX436mLSC1b59b21e"
	baseURL = "https://api.gowinston.ai/v2"
)

func DetectAIContent(text string) (aidetection.AIDetectionResult, error) {
	url := baseURL + "/ai-content-detection"
	payload := map[string]interface{}{
		"text":    text,
		"version": "latest",
	}

	resp, err := makeRequest(url, payload)
	if err != nil {
		return aidetection.AIDetectionResult{}, err
	}

	var result aidetection.AIDetectionResult
	result.OverallAIProbability = float64(resp["score"].(float64))

	if sentences, ok := resp["sentences"].([]interface{}); ok {
		for _, s := range sentences {
			sentence := s.(map[string]interface{})
			if sentence["score"].(float64) > 50 {
				result.FlaggedSections = append(result.FlaggedSections, aidetection.FlaggedSection{
					Text:        sentence["text"].(string),
					Probability: sentence["score"].(float64),
					Reason:      "High AI probability",
				})
			}
		}
	}

	return result, nil
}

func CheckPlagiarism(text string) (plagiarism.PlagiarismResult, error) {
	url := baseURL + "/plagiarism"
	payload := map[string]interface{}{
		"text": text,
	}

	resp, err := makeRequest(url, payload)
	if err != nil {
		return plagiarism.PlagiarismResult{}, err
	}

	var result plagiarism.PlagiarismResult
	result.Score = resp["result"].(map[string]interface{})["score"].(float64)

	if sources, ok := resp["sources"].([]interface{}); ok {
		for _, s := range sources {
			source := s.(map[string]interface{})
			result.Sources = append(result.Sources, plagiarism.Source{
				URL:   source["url"].(string),
				Score: source["score"].(float64),
			})
		}
	}

	return result, nil
}

func makeRequest(url string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
