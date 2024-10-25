package feedbackai

import (
	"context"
	"encoding/json"
	"fmt"
	"psr/cmd/api/secrets"
	"psr/types/feedback"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func GenerateFeedback(personalStatement string) (feedback.FeedbackResponse, error) {
	client := openai.NewClient(secrets.GetEnvVariable("OPEN_AI_SECRET"))

	prompt := fmt.Sprintf(`Analyze the following personal statement and provide detailed feedback on these aspects:
1. Clarity
2. Structure
3. Grammar & Spelling
4. Relevance
5. Engagement
6. Overall Impression

For each aspect, provide a rating out of 10 and detailed feedback. Format your response as a JSON object with the following structure:

{
  "clarity": {"rating": <rating>, "feedback": "<feedback>"},
  "structure": {"rating": <rating>, "feedback": "<feedback>"},
  "grammar_spelling": {"rating": <rating>, "feedback": "<feedback>"},
  "relevance": {"rating": <rating>, "feedback": "<feedback>"},
  "engagement": {"rating": <rating>, "feedback": "<feedback>"},
  "overall_impression": {"rating": <rating>, "feedback": "<feedback>"}
}

Personal Statement:
%s`, personalStatement)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Println("OpenAI API error:", err)
		return feedback.FeedbackResponse{}, fmt.Errorf("OpenAI API error: %v", err)
	}

	jsonResponse := resp.Choices[0].Message.Content

	var feedbackResponse feedback.FeedbackResponse
	err = json.Unmarshal([]byte(jsonResponse), &feedbackResponse)
	if err != nil {
		// If parsing fails, try to extract JSON from the response
		startIndex := strings.Index(jsonResponse, "{")
		endIndex := strings.LastIndex(jsonResponse, "}")
		if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
			jsonResponse = jsonResponse[startIndex : endIndex+1]
			err = json.Unmarshal([]byte(jsonResponse), &feedbackResponse)
			if err != nil {
				return feedback.FeedbackResponse{}, fmt.Errorf("Error parsing extracted JSON response: %v\nExtracted JSON: %s", err, jsonResponse)
			}
		} else {
			return feedback.FeedbackResponse{}, fmt.Errorf("Error parsing JSON response: %v\nRaw response: %s", err, jsonResponse)
		}
	}

	// Validate the parsed response
	if feedbackResponse.Clarity.Rating == 0 || feedbackResponse.Structure.Rating == 0 ||
		feedbackResponse.GrammarSpelling.Rating == 0 || feedbackResponse.Relevance.Rating == 0 ||
		feedbackResponse.Engagement.Rating == 0 || feedbackResponse.OverallImpression.Rating == 0 {
		return feedback.FeedbackResponse{}, fmt.Errorf("Invalid response format from OpenAI API")
	}

	return feedbackResponse, nil
}
