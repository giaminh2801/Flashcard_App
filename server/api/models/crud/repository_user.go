package crud

import (
	"database/sql"
	"time"

	"go-flashcard-api/api/database"
	"go-flashcard-api/api/models"
	"go-flashcard-api/api/security"
	"go-flashcard-api/api/utils/channels"
)

// UserRepository interface for User CRUD
type UserRepository interface {
	Save(models.User) (models.User, error)
	FindAll() ([]models.User, error)
	FindByID(uint64) (models.User, error)
	Update(uint64, models.User) (int64, error)
	Delete(uint64) (int64, error)
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

// FindAll returns all users in our database
func (r *RepositoryUsersCRUD) FindAll() ([]models.User, error) {
	var err error
	users := []models.User{}
	user := models.User{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		rows, err := r.db.Query(database.GetAllUser)
		if err != nil {
			ch <- false
			return
		}
		for rows.Next() {
			err = rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				ch <- false
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return users, nil
	}
	return nil, err
}

// FindByID returns a single user from database
func (r *RepositoryUsersCRUD) FindByID(userID uint64) (models.User, error) {
	var err error
	user := models.User{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		row := r.db.QueryRow(database.GetOneUser, userID)
		err = row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			ch <- false
			return
		}
		if err = row.Err(); err != nil {
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

// Update an user in Database
func (r *RepositoryUsersCRUD) Update(userID uint64, user models.User) (int64, error) {
	var (
		err          error
		rowsAffected int64
	)
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		//Hashing new password
		password, err := security.Hash(user.Password)
		if err != nil {
			ch <- false
			return
		}

		result, err := r.db.Exec(database.UpdateUser, userID, user.Nickname, user.Email, password, time.Now())
		if err != nil {
			ch <- false
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return rowsAffected, nil
	}
	return 0, err
}

// Delete an user from Database
func (r *RepositoryUsersCRUD) Delete(userID uint64) (int64, error) {
	var (
		err          error
		rowsAffected int64
	)
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		result, err := r.db.Exec(database.DeleteUser, userID)
		if err != nil {
			ch <- false
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			ch <- false
			return
		}
		ch <- true

	}(done)
	if channels.OK(done) {
		return rowsAffected, nil
	}
	return 0, err
}
