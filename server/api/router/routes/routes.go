package routes

import (
	"go-flashcard-api/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Route struct
type Route struct {
	URI     string
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

// Load all routes
func Load() []Route {
	routes := make([]Route, 0)
	routes = append(routes, userRoutes...)
	return routes
}

// SetupRoutesWithMiddlewares configurates routes with middlewares
func SetupRoutesWithMiddlewares(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		r.HandleFunc(route.URI,
			middlewares.SetMiddlewareLogger(
				middlewares.SetMiddlewareJSON(
					route.Handler,
				),
			),
		).Methods(route.Method)
	}

	return r
}
