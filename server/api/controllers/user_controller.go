package controllers

import (
	"encoding/json"
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/models/crud"
	"go-flashcard-api/api/responses"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUsers get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)
	func(userRepos crud.UserRepository) {
		users, err := userRepos.FindAll()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.JSON(w, http.StatusOK, users)
	}(repo)
}

// GetUser get an user
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
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
		user, err := userRepos.FindByID(userID)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.JSON(w, http.StatusOK, user)
	}(repo)
}

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

//UpdateUser updates an user in database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
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
		rowsAffected, err := userRepos.Update(userID, user)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		responses.JSON(w, http.StatusOK, rowsAffected)
	}(repo)
}

//DeleteUser deletes an user from database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
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
		rowsAffected, err := userRepos.Delete(userID)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		responses.JSON(w, http.StatusOK, rowsAffected)
	}(repo)
}
