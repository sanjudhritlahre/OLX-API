package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sanjudhritlahre/olx-api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New("file://migrations", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	fmt.Println("Running Migration...")

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		log.Println("Migration Up Completed.")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		log.Println("Migration Down Completed.")
	default:
		log.Fatalf("Unknown Command: %s", os.Args[1])
	}
}
