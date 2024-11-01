package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	aidetection "psr/ai-detection"
	"psr/database/queries"
	"psr/feedbackai"
	plagiarismdetection "psr/plagiarism-detection"
	ftypes "psr/types/feedback"
	statement "psr/types/personal_statement"
	stypes "psr/types/personal_statement"
	"psr/utils/JWT"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/statements/feedback", http.HandlerFunc(h.Feedback)).Methods("POST", "OPTIONS")
	router.HandleFunc("/statements/feedback/{id}", http.HandlerFunc(h.GetFeedback)).Methods("GET")
	router.HandleFunc("/statements/feedback", http.HandlerFunc(h.GetFeedback)).Methods("GET")
}

func (h *Handler) Feedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var request struct {
		Text    string `json:"text"`
		Purpose string `json:"purpose"` // Added purpose field
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Text == "" {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	personalStatement := stypes.PersonalStatement{
		Content: request.Text,
	}

	feedbackResponse, err := feedbackai.GenerateFeedback(personalStatement.Content, request.Purpose) // Pass purpose to feedback generation
	if err != nil {
		fmt.Printf("Error generating feedback: %v\n", err)
		http.Error(w, "Failed to generate feedback", http.StatusInternalServerError)
		return
	}

	aiResult, err := aidetection.DetectAIContent(personalStatement.Content)
	if err != nil {
		fmt.Printf("Error detecting AI content: %v\n", err)
		http.Error(w, "Failed to detect AI content", http.StatusInternalServerError)
		return
	}

	result, err := plagiarismdetection.CheckPlagiarism(personalStatement.Content)
	if err != nil {
		fmt.Printf("Error detecting plagiarism: %v\n", err)
		http.Error(w, "Failed to detect plagiarism", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Plagiarism Percentage: %.2f%%\n", result.PlagiarismPercentage)
	fmt.Println("Sources:")
	for _, source := range result.Sources {
		fmt.Printf("- %s: %.2f%%\n", source.URL, source.Percent)
	}

	// Save AI result
	err = queries.SaveAIResult(1, aiResult)
	if err != nil {
		fmt.Printf("Error saving AI result: %v\n", err)
		http.Error(w, "Failed to save AI result", http.StatusInternalServerError)
		return
	}

	err = queries.SaveStatement(1, personalStatement)
	if err != nil {
		fmt.Printf("Error saving statement: %v\n", err)
		http.Error(w, "Failed to save statement", http.StatusInternalServerError)
		return
	}

	err = queries.SaveFeedback(1, feedbackResponse)
	if err != nil {
		fmt.Printf("Error saving feedback: %v\n", err)
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}
	combinedResponse := ftypes.CombinedResponse{
		Feedback:    feedbackResponse,
		AIDetection: aiResult,
		Plagiarism:  *result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedResponse)
}

type FeedbackResponse struct {
	Statement statement.PersonalStatement `json:"statement"`
	Feedback  ftypes.Feedback             `json:"feedback"`
}

func (h *Handler) GetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	userID, err := JWT.ValidateJWT(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID type", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	feedbackID := vars["id"]

	var feedbacks []ftypes.Feedback
	var statements []statement.PersonalStatement

	if feedbackID != "" {
		// Fetch specific feedback
		feedback, err := queries.GetFeedbackByID(feedbackID)
		if err != nil {
			fmt.Printf("Error retrieving feedback: %v\n", err)
			http.Error(w, "Failed to retrieve feedback", http.StatusInternalServerError)
			return
		}
		feedbacks = append(feedbacks, feedback)

		// Fetch the statement associated with the feedback
		statement, err := queries.GetStatementByFID(feedback.StatementID)
		if err != nil {
			fmt.Printf("Error retrieving statement: %v\n", err)
			http.Error(w, "Failed to retrieve statement", http.StatusInternalServerError)
			return
		}
		statements = append(statements, statement)
	} else {
		// Fetch all feedbacks for the user
		userFeedbacks, err := queries.GetFeedbackForUserStatements(userIDInt)
		if err != nil {
			fmt.Printf("Error retrieving feedbacks: %v\n", err)
			http.Error(w, "Failed to retrieve feedbacks", http.StatusInternalServerError)
			return
		}
		feedbacks = userFeedbacks

		// Fetch all statements for the user
		userStatements, err := queries.GetUserStatements(userIDInt)
		if err != nil {
			fmt.Printf("Error retrieving statements: %v\n", err)
			http.Error(w, "Failed to retrieve statements", http.StatusInternalServerError)
			return
		}
		statements = userStatements
	}

	if len(feedbacks) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No feedbacks found for the user"})
		return
	}

	response := struct {
		Feedbacks  []ftypes.Feedback             `json:"feedbacks"`
		Statements []statement.PersonalStatement `json:"statements"`
	}{
		Feedbacks:  feedbacks,
		Statements: statements,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Feedback ID is required", http.StatusBadRequest)
		return
	}

	err := queries.DeleteFeedback(id)
	if err != nil {
		fmt.Printf("Error deleting feedback: %v\n", err)
		http.Error(w, "Failed to delete feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateStatement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Statement ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Text == "" {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	personalStatement := stypes.PersonalStatement{
		Content: request.Text,
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Statement ID", http.StatusBadRequest)
		return
	}
	err = queries.UpdateStatement(idInt, personalStatement)
	if err != nil {
		fmt.Printf("Error updating statement: %v\n", err)
		http.Error(w, "Failed to update statement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
