package revision

import (
	"encoding/json"
	"net/http"
	"strconv"

	"psr/types/revision"
	ReturnModule "psr/utils/helpful/return_module"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/statements/{id}/revision", http.HandlerFunc(h.Revision)).Methods("POST")
}

func (h *Handler) Revision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statementID := vars["id"]

	// extract feedback from body
	var revision revision.Revision
	err := json.NewDecoder(r.Body).Decode(&revision)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	revision.StatementID, err = strconv.Atoi(statementID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ReturnModule.SendResponse(w, "Feedback", http.StatusOK)
}
