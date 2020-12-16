package routes

import (
	"go-flashcard-api/api/controllers"
	"net/http"
)

var securityRoutes = []Route{
	{
		URI:          "/resetpassword",
		Method:       http.MethodPut,
		Handler:      controllers.ResetPassword,
		AuthRequired: true,
	},
}
