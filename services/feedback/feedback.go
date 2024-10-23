package feedback

import (
	"encoding/json"
	"net/http"
	"strconv"

	"psr/types/feedback"
	ReturnModule "psr/utils/helpful/return_module"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/statements/{id}/feedback", http.HandlerFunc(h.Feedback)).Methods("POST")
}

func (h *Handler) Feedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statementID := vars["id"]

	// extract feedback from body
	var feedback feedback.Feedback
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feedback.StatementID, err = strconv.Atoi(statementID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ReturnModule.SendResponse(w, "Feedback", http.StatusOK)
}
