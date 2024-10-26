package feedback

import (
	aidetection "psr/types/aidetection"
	"time"
)

type Feedback struct {
	ID           int          `json:"id"`
	StatementID  int          `json:"statement_id"`
	FeedbackText FeedbackText `json:"feedback_text"`
	CreatedAt    time.Time    `json:"created_at"`
}

type FeedbackText struct {
	Clarity struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"clarity"`
	Structure struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"structure"`
	GrammarSpelling struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"grammar_spelling"`
	Relevance struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"relevance"`
	Engagement struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"engagement"`
	OverallImpression struct {
		Rating   int    `json:"rating"`
		Feedback string `json:"feedback"`
	} `json:"overall_impression"`
}

type CombinedResponse struct {
	Feedback    FeedbackText                  `json:"feedback"`
	AIDetection aidetection.AIDetectionResult `json:"ai_detection"`
}
