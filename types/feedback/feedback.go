package feedback

import "time"

type Feedback struct {
	ID           int       `json:"id"`
	StatementID  int       `json:"statement_id"`
	FeedbackText string    `json:"feedback_text"`
	CreatedAt    time.Time `json:"created_at"`
}
