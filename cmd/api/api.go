package api

import (
	"fmt"
	"net/http"

	"psr/database"
	routes "psr/services"
	openai "psr/utils/ai/openai"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	database.InitializeDatabase()
	if !database.CreateTables() {
		fmt.Println("Error creating tables")
		return fmt.Errorf("Error creating tables")
	}
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/").Subrouter()

	routeHandler := routes.NewRouteHandler()
	routeHandler.RegisterRoutes(subrouter)
	go openai.Init()
	fmt.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)

}
