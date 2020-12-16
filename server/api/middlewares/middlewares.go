package middlewares

import (
	"context"
	"go-flashcard-api/api/auth"
	"go-flashcard-api/api/controllers"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/responses"
	"go-flashcard-api/api/utils/channels"
	"go-flashcard-api/api/utils/types"
	"log"
	"net/http"
)

// SetMiddlewareLogger displays a info message of the API
func SetMiddlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)
		next(w, r)
	}
}

// SetMiddlewareJSON set the application Content-Type
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			responses.JSON(w, http.StatusNoContent, nil)
			return
		}

		next(w, r)
	}
}

// SetMiddlewareAuthentication authorize an access token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := auth.VerifyAccessToken(w, r)
		if accessToken == nil {
			refreshToken := auth.VerifyRefreshToken(w, r)
			if refreshToken == nil {
				return
			}
			done := make(chan bool)
			var err error
			var ctx context.Context
			go controllers.Refresh(done, &err, &ctx, refreshToken, w, r)
			if err != nil {
				responses.ERROR(w, http.StatusUnauthorized, err)
			}
			if channels.OK(done) {
				next(w, r.WithContext(ctx))
			}
			return
		} else if accessToken.Valid {
			ctx := context.WithValue(
				r.Context(),
				types.StringKey("user"),
				accessToken.Claims.(*models.Claim).User,
			)
			next(w, r.WithContext(ctx))
		}
	}
}
