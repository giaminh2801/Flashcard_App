package routes

import (
	"go-flashcard-api/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Route struct
type Route struct {
	URI          string
	Method       string
	Handler      func(w http.ResponseWriter, r *http.Request)
	AuthRequired bool
}

// Load all routes
func Load() []Route {
	routes := make([]Route, 0)
	routes = append(routes, userRoutes...)
	routes = append(routes, loginRoutes...)
	return routes
}

// SetupRoutesWithMiddlewares configurates routes with middlewares
func SetupRoutesWithMiddlewares(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		if route.AuthRequired {
			r.HandleFunc(route.URI,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(
						middlewares.SetMiddlewareAuthentication(route.Handler),
					),
				),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(route.Handler),
				),
			).Methods(route.Method)
		}
	}

	return r
}
