//go:build integration

package testutil

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikpopkov/running-club/api/internal/pkg/migrate"
	"github.com/stretchr/testify/require"
)

var (
	migrateOnce sync.Once
	migrateErr  error
)

func Open(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pulse:pulse@localhost:5432/running_club_test?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Skipf("no db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Skipf("db unavailable: %v", err)
	}
	_ = db.Close()

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	Migrate(t, dsn)
	Truncate(t, pool)
	return pool
}

func Migrate(t *testing.T, dsn string) {
	t.Helper()
	migrateOnce.Do(func() {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			migrateErr = err
			return
		}
		defer db.Close()
		_, _ = db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
		migrateErr = migrate.Up(db, migrationsDir())
	})
	require.NoError(t, migrateErr)
}

func Truncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE
			activity_streams,
			activities,
			segments,
			workouts,
			announce_signups,
			announces,
			prs,
			races,
			month_stats,
			plan_weeks,
			user_integrations,
			memberships,
			clubs,
			users
		RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
}

func migrationsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "scripts", "migrations")
}
