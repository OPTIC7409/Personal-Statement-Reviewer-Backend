package user

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
	router.HandleFunc("/users/me", http.HandlerFunc(h.GetUser)).Methods("GET")
	router.HandleFunc("/users/me", http.HandlerFunc(h.UpdateUser)).Methods("PUT")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ReturnModule.SendResponse(w, "User", http.StatusOK)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	ReturnModule.SendResponse(w, "Statement", http.StatusOK)
}
