package controllers

import (
	"encoding/json"
	"go-flashcard-api/api/auth"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/responses"
	"net/http"
)

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := auth.SignIn(user.Email, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, token)
}
