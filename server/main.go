package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"y/handler"
	"y/router"
)

func main() {
	// Eagerly initialize the DB pool so we fail fast if Postgres is unreachable.
	handler.DB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	r := router.Router()
	fmt.Printf("Starting server on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
