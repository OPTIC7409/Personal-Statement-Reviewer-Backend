package routes

import (
	"psr/services/auth"
	"psr/services/dashboard"
	"psr/services/feedback"
	"psr/services/revision"
	"psr/services/statements"
	"psr/services/user"

	"github.com/gorilla/mux"
)

type RouteHandler struct{}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (rh *RouteHandler) RegisterRoutes(router *mux.Router) {
	authHandler := auth.NewHandler()
	authHandler.RegisterRoutes(router)

	dashboardHandler := dashboard.NewHandler()
	dashboardHandler.RegisterRoutes(router)

	feedbackHandler := feedback.NewHandler()
	feedbackHandler.RegisterRoutes(router)

	revisionHandler := revision.NewHandler()
	revisionHandler.RegisterRoutes(router)

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(router)

	statementHandler := statements.NewHandler()
	statementHandler.RegisterRoutes(router)

}
