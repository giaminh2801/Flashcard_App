package routes

import (
	controlers "go-flashcard-api/api/controllers"
	"net/http"
)

var loginRoutes = []Route{
	{
		URI:          "/login",
		Method:       http.MethodPost,
		Handler:      controlers.Login,
		AuthRequired: false,
	},
}
