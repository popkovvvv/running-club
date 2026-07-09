package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikpopkov/running-club/api/internal/config"
	"github.com/nikpopkov/running-club/api/internal/pkg/migrate"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dir := "scripts/migrations"
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|down]")
	}
	switch os.Args[1] {
	case "up":
		if err := migrate.Up(db, dir); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := migrate.Down(db, dir); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("usage: migrate [up|down]")
	}
}
