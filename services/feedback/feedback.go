package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"

	"psr/database/queries"
	"psr/feedbackai"
	statement "psr/types/personal_statement"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/statements/feedback", http.HandlerFunc(h.Feedback)).Methods("POST", "OPTIONS")
}

func (h *Handler) Feedback(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
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

	personalStatement := statement.PersonalStatement{
		Content: request.Text,
	}

	feedbackResponse, err := feedbackai.GenerateFeedback(personalStatement.Content)
	if err != nil {
		fmt.Printf("Error generating feedback: %v\n", err)
		http.Error(w, "Failed to generate feedback", http.StatusInternalServerError)
		return
	}

	err = queries.SaveStatement(personalStatement)
	if err != nil {
		fmt.Printf("Error saving statement: %v\n", err)
		http.Error(w, "Failed to save statement", http.StatusInternalServerError)
		return
	}

	err = queries.SaveFeedback(feedbackResponse)
	if err != nil {
		fmt.Printf("Error saving feedback: %v\n", err)
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbackResponse)
}
