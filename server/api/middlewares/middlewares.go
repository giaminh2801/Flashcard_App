package middlewares

import (
	"context"
	"go-flashcard-api/api/auth"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/utils/types"
	"go-flashcard-api/config"
	"log"
	"net/http"
)

// SetMiddlewareLogger displays a info message of the API
func SetMiddlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)
		next(w, r)
	}
}

// SetMiddlewareJSON set the application Content-Type
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.CLIENTURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication authorize an access
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := auth.ExtractToken(w, r)
		if token == nil {
			return
		}
		if token.Valid {
			ctx := context.WithValue(
				r.Context(),
				types.UserKey("user"),
				token.Claims.(*models.Claim).User,
			)
			next(w, r.WithContext(ctx))
		}
	}
}
