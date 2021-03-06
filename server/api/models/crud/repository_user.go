package crud

import (
	"database/sql"
	"errors"
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
	UpdatePassword(uint64, string, string) (int64, error)
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
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		if err = user.BeforeSave(); err != nil {
			*p_err = err
			ch <- false
			return
		}
		_, err = r.db.Exec(database.CreateUserQuery, user.Nickname, user.Email, user.Password)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
				err = errors.New("Email is already existed, please use another email")
			}
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)
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
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		rows, err := r.db.Query(database.GetAllUser)
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		for rows.Next() {
			err = rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				*p_err = err
				ch <- false
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)
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
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		row := r.db.QueryRow(database.GetOneUser, userID)
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
		ch <- true
	}(done, &err)
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
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		result, err := r.db.Exec(database.UpdateUser, userID, user.Nickname, user.Email, time.Now())
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
				err = errors.New("Email is already existed, please use another email")
			}
			*p_err = err
			ch <- false
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)
	if channels.OK(done) {
		return rowsAffected, nil
	}
	return 0, err
}

// UpdatePassword updates user's password in Database
func (r *RepositoryUsersCRUD) UpdatePassword(userID uint64, oldPassword string, newPassword string) (int64, error) {
	var (
		err               error
		rowsAffected      int64
		oldHashedPassword string
	)
	done := make(chan bool)
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		row := r.db.QueryRow(database.GetPassword, userID)
		err = row.Scan(&oldHashedPassword)
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
		err = security.VerifyPassword(oldHashedPassword, oldPassword)
		if err != nil {
			err = errors.New("Your old password isn't correct")
			*p_err = err
			ch <- false
			return
		}
		newHashedPassword, err := security.Hash(newPassword)
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		result, err := r.db.Exec(database.UpdatePassword, userID, newHashedPassword, time.Now())
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)
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
	go func(ch chan<- bool, p_err *error) {
		defer close(ch)
		result, err := r.db.Exec(database.DeleteUser, userID)
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			*p_err = err
			ch <- false
			return
		}
		ch <- true
	}(done, &err)
	if channels.OK(done) {
		return rowsAffected, nil
	}
	return 0, err
}
