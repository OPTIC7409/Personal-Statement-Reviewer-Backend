package dashboard

import (
	"net/http"

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

	ReturnModule.SendResponse(w, "Dashboard", http.StatusOK)
}
