package stripeendpoint

import (
	"psr/stripe"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/stripe/webhook", stripe.HandleStripeWebhook).Methods("POST")
}
