package api

import (
	"fmt"
	"net/http"
	"strings"

	"psr/database"
	routes "psr/services"
	"psr/stripe"
	openai "psr/utils/ai/openai"
	utils "psr/utils/supertokens"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	database.InitializeDatabase()
	if !database.CreateTables() {
		fmt.Println("Error creating tables")
		return fmt.Errorf("Error creating tables")
	}

	err := utils.InitializeSuperTokens("http://localhost:3002", "http://localhost:3000", database.GetConnection())
	if err != nil {
		return fmt.Errorf("Failed to initialize SuperTokens: %v", err)
	}

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api").Subrouter()

	routeHandler := routes.NewRouteHandler()
	routeHandler.RegisterRoutes(subrouter)

	go openai.Init()
	go stripe.InitStripe()
	fmt.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, corsMiddleware(
		supertokens.Middleware(router)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			strings.Join(append([]string{"Content-Type", "Authorization"},
				supertokens.GetAllCORSHeaders()...), ","))

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
