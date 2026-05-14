package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Tiny helper: prints a bcrypt hash for each password passed as a CLI arg.
// Used once to generate the seed hashes baked into db/init.sql.
// Run with: go run ./cmd/hashgen admin123 demo123
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: hashgen <password> [<password>...]")
		os.Exit(1)
	}
	for _, p := range os.Args[1:] {
		h, err := bcrypt.GenerateFromPassword([]byte(p), 10)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("%s -> %s\n", p, string(h))
	}
}
