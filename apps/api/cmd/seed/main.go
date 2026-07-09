package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/config"
	"github.com/nikpopkov/running-club/api/internal/pkg/seed"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := seed.Run(ctx, pool); err != nil {
		log.Fatal(err)
	}
	log.Println("seed ok")
}
