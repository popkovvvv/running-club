package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Up(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}
	return nil
}

func Down(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}
	if err := goose.Down(db, dir); err != nil {
		return fmt.Errorf("goose.Down: %w", err)
	}
	return nil
}
