package controllers

import (
	"encoding/json"
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/models/crud"
	"go-flashcard-api/api/responses"
	"go-flashcard-api/api/utils/types"
	"net/http"
)

type resetPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ResetPassword handler
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	resetPassword := resetPassword{}
	err := json.NewDecoder(r.Body).Decode(&resetPassword)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)
	userID := r.Context().Value(types.StringKey("user")).(models.User).ID
	func(userRepos crud.UserRepository) {
		rowsAffected, err := userRepos.UpdatePassword(userID, resetPassword.OldPassword, resetPassword.NewPassword)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		responses.JSON(w, http.StatusOK, rowsAffected)
	}(repo)
}
