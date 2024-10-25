package feedbackai

import (
	"context"
	"encoding/json"
	"fmt"
	"psr/types/feedback"
	"psr/utils/helpful/parsing"

	"github.com/sashabaranov/go-openai"
)

func GenerateFeedback(personalStatement string) (feedback.FeedbackResponse, error) {
	client := openai.NewClient("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")

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
		return feedback.FeedbackResponse{}, fmt.Errorf("OpenAI API error: %v", err)
	}

	jsonResponse := resp.Choices[0].Message.Content

	var feedbackResponse feedback.FeedbackResponse
	err = json.Unmarshal([]byte(jsonResponse), &feedbackResponse)
	if err != nil {
		// If parsing fails, try to extract JSON from the response
		err = parsing.ExtractJSONToStruct(jsonResponse, &feedbackResponse)
		if err != nil {
			return feedback.FeedbackResponse{}, fmt.Errorf("Error extracting JSON: %v\nRaw response: %s", err, jsonResponse)
		}
	}

	// Validate the parsed response
	if !isValidFeedbackResponse(feedbackResponse) {
		return feedback.FeedbackResponse{}, fmt.Errorf("Invalid response format from OpenAI API")
	}

	return feedbackResponse, nil
}

func isValidFeedbackResponse(fr feedback.FeedbackResponse) bool {
	return fr.Clarity.Rating != 0 &&
		fr.Structure.Rating != 0 &&
		fr.GrammarSpelling.Rating != 0 &&
		fr.Relevance.Rating != 0 &&
		fr.Engagement.Rating != 0 &&
		fr.OverallImpression.Rating != 0
}
