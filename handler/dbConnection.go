package handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// package handler

// import (
//     "database/sql"
//     _ "github.com/lib/pq" // Import the PostgreSQL driver
//     "log"
// )

// func createConnection() *sql.DB {
//     connStr := "postgres://postgres:mysecretpassword@localhost:5432/codedb?sslmode=disable"
//     db, err := sql.Open("postgres", connStr)
//     if err != nil {
//         log.Fatalf("Error opening database: %v", err)
//     }
//     return db
// }
