package routes

import (
	"github.com/gorilla/mux"
)

type RouteHandler struct{}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (rh *RouteHandler) RegisterRoutes(router *mux.Router) {

}
