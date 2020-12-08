package routes

import (
	controlers "go-flashcard-api/api/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Handler:      controlers.GetUsers,
		AuthRequired: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Handler:      controlers.CreateUser,
		AuthRequired: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Handler:      controlers.GetUser,
		AuthRequired: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Handler:      controlers.UpdateUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Handler:      controlers.DeleteUser,
		AuthRequired: true,
	},
}
