package database

import (
	"database/sql"
	"log"
)

//CreateUsersTable creates users table
func CreateUsersTable(DB *sql.DB) {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users( 
			id BIGSERIAL PRIMARY KEY, 
			nickname VARCHAR (50) NOT NULL, 
			email VARCHAR (50) UNIQUE NOT NULL, 
			password VARCHAR (300) NOT NULL,
			created_on TIMESTAMP NOT NULL default current_timestamp,
			updated_at TIMESTAMP NOT NULL default current_timestamp)`)
	if err != nil {
		log.Fatal(err)
		return
	}
}
