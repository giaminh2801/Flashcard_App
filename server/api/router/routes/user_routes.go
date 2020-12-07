package routes

import (
	controlers "go-flashcard-api/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:     "/users",
		Method:  http.MethodPost,
		Handler: controlers.CreateUser,
	},
	{
		URI:     "/users/{id}",
		Method:  http.MethodGet,
		Handler: controlers.GetUser,
	},
}
