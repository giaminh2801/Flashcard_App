package crud

import (
	"database/sql"
	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/utils/channels"
)

// UserRepository interface for User CRUD
type UserRepository interface {
	Save(models.User) (models.User, error)
}

// RepositoryUsersCRUD is the struct for the User CRUD
type RepositoryUsersCRUD struct {
	db *sql.DB
}

// NewRepositoryUsersCRUD returns a new repository with DB connection
func NewRepositoryUsersCRUD(db *sql.DB) *RepositoryUsersCRUD {
	return &RepositoryUsersCRUD{db}
}

// Save returns a new user created or an error
func (r *RepositoryUsersCRUD) Save(user models.User) (models.User, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		if err = user.BeforeSave(); err != nil {
			ch <- false
			return
		}
		_, err = r.db.Exec(database.CreateUserQuery, user.Nickname, user.Email, user.Password)
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return user, nil
	}
	return models.User{}, err
}
