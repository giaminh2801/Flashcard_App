package routes

import (
	"go-flashcard-api/api/controllers"
	"net/http"
)

var authRoutes = []Route{
	{
		URI:          "/login",
		Method:       http.MethodPost,
		Handler:      controllers.Login,
		AuthRequired: false,
	},
	{
		URI:          "/logout",
		Method:       http.MethodDelete,
		Handler:      controllers.Logout,
		AuthRequired: true,
	},
}
