package queries

import (
	"encoding/json"
	"fmt"
	"psr/database"
	"psr/types/aidetection"
	"psr/types/feedback"
)

func GetFeedbackForUserStatements(userID int) ([]feedback.Feedback, error) {
	var feedbacks []feedback.Feedback
	rows, err := database.GetConnection().Query(`
		SELECT f.id, f.statement_id, f.feedback_text, f.created_at
		FROM feedback f
		JOIN personal_statements ps ON f.statement_id = ps.id
		WHERE ps.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("Error querying feedback: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f feedback.Feedback
		var feedbackText string
		err := rows.Scan(&f.ID, &f.StatementID, &feedbackText, &f.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error scanning feedback row: %v", err)
		}
		err = json.Unmarshal([]byte(feedbackText), &f.FeedbackText)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling feedback JSON: %v", err)
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func SaveFeedback(statementID int, feedback feedback.FeedbackText) error {
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("Error marshaling feedback: %v", err)
	}

	_, err = database.GetConnection().Exec(`
		INSERT INTO feedback (statement_id, feedback_text)
		VALUES ($1, $2)
	`, statementID, string(feedbackJSON))

	if err != nil {
		return fmt.Errorf("Error saving feedback: %v", err)
	}

	return nil
}

func SaveAIResult(statementID int, result aidetection.AIDetectionResult) error {
	flaggedSectionsJSON, err := json.Marshal(result.FlaggedSections)
	if err != nil {
		return fmt.Errorf("Error marshaling flagged sections: %v", err)
	}

	_, err = database.GetConnection().Exec(`
		INSERT INTO ai_detection (statement_id, overall_ai_probability, flagged_sections)
		VALUES ($1, $2, $3)
	`, statementID, result.OverallAIProbability, string(flaggedSectionsJSON))

	if err != nil {
		return fmt.Errorf("Error saving AI result: %v", err)
	}

	return nil
}

func GetFeedbackByID(id string) (feedback.Feedback, error) {
	var f feedback.Feedback
	var feedbackText string
	err := database.GetConnection().QueryRow(`
		SELECT id, statement_id, feedback_text, created_at
		FROM feedback
		WHERE id = $1
	`, id).Scan(&f.ID, &f.StatementID, &feedbackText, &f.CreatedAt)
	if err != nil {
		return feedback.Feedback{}, fmt.Errorf("Error querying feedback: %v", err)
	}

	err = json.Unmarshal([]byte(feedbackText), &f.FeedbackText)
	if err != nil {
		return feedback.Feedback{}, fmt.Errorf("Error unmarshaling feedback JSON: %v", err)
	}

	return f, nil
}

func GetFeedbackBySID(statementID int) (feedback.Feedback, error) {
	var f feedback.Feedback
	err := database.GetConnection().QueryRow(`
		SELECT id, statement_id, feedback_text, created_at
		FROM feedback
		WHERE statement_id = $1
	`, statementID).Scan(&f.ID, &f.StatementID, &f.FeedbackText, &f.CreatedAt)
	if err != nil {
		return feedback.Feedback{}, fmt.Errorf("Error querying feedback: %v", err)
	}

	return f, nil
}

func DeleteFeedback(id string) error {
	_, err := database.GetConnection().Exec(`DELETE FROM feedback WHERE id = $1`, id)
	return err
}
