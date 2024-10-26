package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	aidetection "psr/ai-detection"
	"psr/database/queries"
	"psr/feedbackai"
	ftypes "psr/types/feedback"
	statement "psr/types/personal_statement"
	stypes "psr/types/personal_statement"

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

	feedbackResponse, err := feedbackai.GenerateFeedback(personalStatement.Content)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		userID := 1
		statements, err := queries.GetUserStatements(userID)
		if err != nil {
			fmt.Printf("Error retrieving statements: %v\n", err)
			http.Error(w, "Failed to retrieve statements", http.StatusInternalServerError)
			return
		}

		if len(statements) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"user does not have any": ""})
			return
		}

		var combinedResponses []FeedbackResponse
		for _, statement := range statements {
			feedbackResponse, err := queries.GetFeedbackBySID(statement.ID)
			if err != nil {
				fmt.Printf("Error retrieving feedback for statement ID %d: %v\n", statement.ID, err)
				continue
			}

			combinedResponses = append(combinedResponses, FeedbackResponse{
				Statement: statement,
				Feedback:  feedbackResponse,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedResponses)
		return
	} else {
		feedbackResponse, err := queries.GetFeedbackByID(id)
		if err != nil {
			fmt.Printf("Error retrieving feedback: %v\n", err)
			http.Error(w, "Failed to retrieve feedback", http.StatusInternalServerError)
			return
		}

		statementID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error converting feedback ID to int: %v\n", err)
			http.Error(w, "Failed to convert feedback ID to int", http.StatusBadRequest)
			return
		}
		statement, err := queries.GetStatementByFID(statementID)
		if err != nil {
			fmt.Printf("Error retrieving statement for feedback ID %s: %v\n", id, err)
			http.Error(w, "Failed to retrieve statement", http.StatusInternalServerError)
			return
		}

		combinedResponses := []FeedbackResponse{{
			Statement: statement,
			Feedback:  feedbackResponse,
		}}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedResponses)
	}
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
