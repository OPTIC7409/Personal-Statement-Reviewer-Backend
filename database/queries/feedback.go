package queries

import (
	"encoding/json"
	"fmt"
	"psr/database"
	"psr/types/feedback"
)

func GetFeedbackForUserStatements(userID int) ([]feedback.Feedback, error) {
	var feedbacks []feedback.Feedback
	_, err := database.GetConnection().Exec(`
		SELECT * FROM feedback WHERE user_id = $1
	`, userID)
	return feedbacks, err
}

func SaveFeedback(feedback feedback.FeedbackResponse) error {
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("Error marshaling feedback: %v", err)
	}

	_, err = database.GetConnection().Exec(`
		INSERT INTO feedback (feedback_text)
		VALUES (?)
	`, string(feedbackJSON))

	if err != nil {
		return fmt.Errorf("Error saving feedback: %v", err)
	}

	return nil
}
