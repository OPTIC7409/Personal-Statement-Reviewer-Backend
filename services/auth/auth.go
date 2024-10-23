package auth

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
	router.HandleFunc("/auth/signup", http.HandlerFunc(h.Signup)).Methods("GET")
	router.Handle("/auth/login", http.HandlerFunc(h.Login)).Methods("POST")
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {

	ReturnModule.SendResponse(w, "Signup successful", http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	ReturnModule.SendResponse(w, "Login successful", http.StatusOK)
}
