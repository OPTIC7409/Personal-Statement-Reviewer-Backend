package api

import (
	"fmt"
	"net/http"
	"strings"

	"psr/database"
	routes "psr/services"
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
		fmt.Errorf("Failed to initialize SuperTokens: %v", err)
	}

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api").Subrouter()

	routeHandler := routes.NewRouteHandler()
	routeHandler.RegisterRoutes(subrouter)

	go openai.Init()
	fmt.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, corsMiddleware(
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			router.ServeHTTP(rw, r)
		}))))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "https://localhost:3000")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			// we add content-type + other headers used by SuperTokens
			response.Header().Set("Access-Control-Allow-Headers",
				strings.Join(append([]string{"Content-Type"},
					supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}
