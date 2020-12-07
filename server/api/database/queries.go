package database

// QUERIES
const (
	QUERIES         = ""
	CheckUserExists = `SELECT id from users WHERE email = $1`
	CreateUserQuery = `INSERT INTO users(id, nickname, email, password) VALUES (DEFAULT, $1, $2, $3)`
)
