package statements

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
	router.HandleFunc("/statements/{id}", http.HandlerFunc(h.GetStatement)).Methods("GET")
	router.HandleFunc("/statements/", http.HandlerFunc(h.NewStatement)).Methods("POST")
}

func (h *Handler) GetStatement(w http.ResponseWriter, r *http.Request) {
	ReturnModule.SendResponse(w, "Statement", http.StatusOK)
}

func (h *Handler) NewStatement(w http.ResponseWriter, r *http.Request) {

	ReturnModule.SendResponse(w, "Statement", http.StatusOK)
}
