package database

// QUERIES users
const (
	QUERIES         = ""
	CheckUserExists = `SELECT id FROM users WHERE email = $1`
	CreateUserQuery = `INSERT INTO users(id, nickname, email, password) VALUES (DEFAULT, $1, $2, $3)`
	LoginQuery      = `SELECT * FROM users WHERE email = $1`
	GetAllUser      = `SELECT * FROM users`
	GetOneUser      = `SELECT * FROM users WHERE id = $1`
	GetPassword     = `SELECT password FROM users WHERE id = $1`
	UpdateUser      = `UPDATE users SET nickname=$2, email=$3, updated_at=$4 WHERE id=$1`
	UpdatePassword  = `UPDATE users SET password=$2, updated_at=$3 WHERE id=$1`
	DeleteUser      = `DELETE FROM users WHERE id = $1`
)
