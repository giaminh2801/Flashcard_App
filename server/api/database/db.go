package database

import (
	"database/sql"
	"go-flashcard-api/config"

	_ "github.com/lib/pq" // postgresql driver
)

// Connect to database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		return nil, err
	}

	CreateUsersTable(db)

	return db, nil
}
