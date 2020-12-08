package auth

import (
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/security"
	"go-flashcard-api/api/utils/channels"
)

// SignIn method
func SignIn(email, password string) (string, error) {
	user := models.User{}
	var err error
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		db, err := database.Connect()
		if err != nil {
			ch <- false
			return
		}
		defer db.Close()

		row := db.QueryRow(database.LoginQuery, email)
		err = row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password)
		if err != nil {
			ch <- false
			return
		}
		if err = row.Err(); err != nil {
			ch <- false
			return
		}

		err = security.VerifyPassword(user.Password, password)
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		user.Password = ""
		return GenerateJWT(user)
	}

	return "", err
}
