package controlers

import (
	"encoding/json"
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/models/crud"
	"go-flashcard-api/api/responses"
	"net/http"
)

// CreateUser inserts new User to DB if valid
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
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

	func(userRepos crud.UserRepository) {
		user, err = userRepos.Save(user)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.JSON(w, http.StatusCreated, user)
	}(repo)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	return
}
