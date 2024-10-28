package stripe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

func InitStripe() {
	stripe.Key = "sk_test_51QEyq6B3f5MlzO9LK0oZI4YJ9SqBse7NTKYSRHCw5fsOHXLWCl0Ni1ntof8D0F18sg8sbaImZsBmTTw51CDHUtkd00TfpArnQa"

}

func HandleStripeWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := ""
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handleCheckoutSessionCompleted(session)
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handleInvoicePaymentSucceeded(invoice)
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handleSubscriptionDeleted(subscription)
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCheckoutSessionCompleted(session stripe.CheckoutSession) {
	// Handle successful checkout
	// Update user's subscription status in your database
	// Provision access to the purchased plan
	fmt.Printf("Checkout completed for session %s\n", session.ID)
}

func handleInvoicePaymentSucceeded(invoice stripe.Invoice) {
	// Handle successful invoice payment
	// Update payment status in your database
	fmt.Printf("Invoice payment succeeded for invoice %s\n", invoice.ID)
}

func handleSubscriptionDeleted(subscription stripe.Subscription) {
	// Handle subscription deletion
	// Update subscription status in your database
	fmt.Printf("Subscription deleted for subscription %s\n", subscription.ID)
}
