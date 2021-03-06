package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// PORT server port
var (
	PORT      = 0
	ATSECRET  []byte // Access token secret key
	RTSECRET  []byte // Refresh token secret key
	DBURL     = ""
	CLIENTURL = "http://localhost:3000"
)

// Load configurations
func Load() {
	var err error
	err = godotenv.Load("config_var.env")
	if err != nil {
		log.Fatal(err)
	}
	PORT, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		PORT = 8000
	}
	DBURL = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	ATSECRET = []byte(os.Getenv("ACCESS_SECRET"))
	RTSECRET = []byte(os.Getenv("REFRESH_SECRET"))
}
