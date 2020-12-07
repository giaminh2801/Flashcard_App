package api

import (
	"fmt"
	"go-flashcard-api/api/router"
	"go-flashcard-api/config"
	"log"
	"net/http"
)

func init() {
	config.Load()
}

func Run() {
	fmt.Printf("\n\tListening [::]:%d", config.PORT)
	listen(config.PORT)
}

func listen(port int) {
	r := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
