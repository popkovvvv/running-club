package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikpopkov/running-club/api/cmd/api/service_provider"
	"github.com/nikpopkov/running-club/api/internal/config"
	"github.com/nikpopkov/running-club/api/internal/pkg/migrate"
	"github.com/nikpopkov/running-club/api/internal/pkg/seed"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	sp := service_provider.New(cfg)
	if err := sp.Boot(ctx); err != nil {
		log.Fatalf("boot: %v", err)
	}
	defer sp.Close()

	if cfg.Seed {
		if err := seed.Run(ctx, sp.Pool()); err != nil {
			log.Fatalf("seed: %v", err)
		}
	}

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: sp.Handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("api listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
}

func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	dir := migrationsDir()
	return migrate.Up(db, dir)
}

func migrationsDir() string {
	candidates := []string{
		"scripts/migrations",
		"apps/api/scripts/migrations",
		filepath.Join(os.Getenv("PWD"), "scripts/migrations"),
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}
	return "scripts/migrations"
}
