package dashboard

import (
	"net/http"

	"psr/database/queries"
	ffeedback "psr/types/feedback"
	statement "psr/types/personal_statement"
	ReturnModule "psr/utils/helpful/return_module"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dashboard", http.HandlerFunc(h.Dashboard)).Methods("GET")
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(int) // Assuming user ID is stored in the request context

	statements, err := queries.GetUserStatements(userID)
	if err != nil {
		ReturnModule.SendResponse(w, "Failed to retrieve user statements", http.StatusInternalServerError)
		return
	}

	feedback, err := queries.GetFeedbackForUserStatements(userID)
	if err != nil {
		ReturnModule.SendResponse(w, "Failed to retrieve feedback for user statements", http.StatusInternalServerError)
		return
	}

	dashboardData := struct {
		Statements []statement.PersonalStatement `json:"statements"`
		Feedback   []ffeedback.Feedback          `json:"feedback"`
	}{
		Statements: statements,
		Feedback:   feedback,
	}

	ReturnModule.SendResponse(w, dashboardData, http.StatusOK)
}
