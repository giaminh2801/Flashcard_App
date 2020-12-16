package router

import (
	"go-flashcard-api/api/router/routes"
	"go-flashcard-api/config"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// w.Header().Set("Access-Control-Allow-Origin", config.CLIENTURL)
// w.Header().Set("Access-Control-Allow-Credentials", "true")
// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

// New routes
func New() (*mux.Router, func(http.Handler) http.Handler) {
	r := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"})
	origins := handlers.AllowedOrigins([]string{config.CLIENTURL})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})
	handler := handlers.CORS(headers, origins, methods)
	return routes.SetupRoutesWithMiddlewares(r), handler
}
