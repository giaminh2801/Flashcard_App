package auth

import (
	"errors"
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/security"
	"go-flashcard-api/api/utils/channels"
)

// SignIn method
func SignIn(email, password string) (models.User, *TokenDetails, error) {
	user := models.User{}
	var err error
	done := make(chan bool)

	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		db, err := database.Connect()
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		defer db.Close()

		row := db.QueryRow(database.LoginQuery, email)
		err = row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		if err = row.Err(); err != nil {
			*p_err = err
			ch <- false
			return
		}

		err = security.VerifyPassword(user.Password, password)
		if err != nil {
			err = errors.New("Password isn't correct, please try again")
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)

	if channels.OK(done) {
		user.Password = ""
		tokenDetails, err := GenerateJWT(user)
		return user, tokenDetails, err
	}
	return models.User{}, nil, err
}
