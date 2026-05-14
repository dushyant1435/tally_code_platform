package handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	dbOnce     sync.Once
)

// DB returns a process-wide *sql.DB pool. The pool is initialized lazily on the
// first call and reused for the lifetime of the server.
func DB() *sql.DB {
	dbOnce.Do(func() {
		// Load .env if present; ignore the error so the server still starts when
		// running inside Docker where env vars are injected directly.
		_ = godotenv.Load(".env")

		url := os.Getenv("POSTGRES_URL")
		if url == "" {
			log.Fatal("POSTGRES_URL is not set")
		}

		db, err := sql.Open("postgres", url)
		if err != nil {
			log.Fatalf("failed to open postgres connection: %v", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Retry the initial ping so we wait for Postgres in docker-compose.
		var pingErr error
		for i := 0; i < 30; i++ {
			pingErr = db.Ping()
			if pingErr == nil {
				break
			}
			log.Printf("waiting for postgres (%d/30): %v", i+1, pingErr)
			time.Sleep(time.Second)
		}
		if pingErr != nil {
			log.Fatalf("postgres unreachable: %v", pingErr)
		}

		fmt.Println("Successfully connected to Postgres")
		dbInstance = db
	})
	return dbInstance
}
