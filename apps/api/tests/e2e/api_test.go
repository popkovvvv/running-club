//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nikpopkov/running-club/api/cmd/api/service_provider"
	"github.com/nikpopkov/running-club/api/internal/config"
	"github.com/nikpopkov/running-club/api/internal/pkg/migrate"
	"github.com/nikpopkov/running-club/api/internal/pkg/seed"
	"github.com/stretchr/testify/require"
)

func TestAPIFlows(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pulse:pulse@localhost:5432/running_club_test?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Skipf("no db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Skipf("db unavailable: %v", err)
	}
	defer db.Close()

	_, _ = db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	migDir := filepath.Join("..", "..", "scripts", "migrations")
	require.NoError(t, migrate.Up(db, migDir))

	cfg := config.Config{DatabaseURL: dsn, JWTSecret: "test-secret", HTTPAddr: ":0", Seed: true}
	sp := service_provider.New(cfg)
	require.NoError(t, sp.Boot(context.Background()))
	defer sp.Close()
	require.NoError(t, seed.Run(context.Background(), sp.Pool()))
	h := sp.Handler()

	coachTok := login(t, h, "coach@pulse.run", "password")
	athTok := login(t, h, "nikita@pulse.run", "password")

	t.Run("signup_cta", func(t *testing.T) {
		ann := getJSON[[]map[string]any](t, h, "GET", "/api/v1/announces", athTok, nil)
		require.NotEmpty(t, ann)
		id := ann[0]["id"].(string)
		require.Equal(t, "Записаться", ann[0]["scheduleCta"])
		signed := getJSON[map[string]any](t, h, "POST", "/api/v1/announces/"+id+"/signup", athTok, nil)
		require.Equal(t, "Вы записаны", signed["scheduleCta"])
	})

	t.Run("calendar_blank", func(t *testing.T) {
		cal := getJSON[map[string]any](t, h, "GET", "/api/v1/schedule/calendar", athTok, nil)
		cells := cal["cells"].([]any)
		foundBlank := false
		for _, c := range cells {
			m := c.(map[string]any)
			if m["blank"] == true {
				foundBlank = true
				require.Equal(t, "transparent", m["dot"])
				require.Equal(t, "transparent", m["bg"])
			}
		}
		require.True(t, foundBlank)
	})

	t.Run("coach_palette_and_remove", func(t *testing.T) {
		club := getJSON[map[string]any](t, h, "PATCH", "/api/v1/club/palette", coachTok, map[string]string{"accentHex": "#c8ff34"})
		require.Equal(t, "#c8ff34", club["accentHex"])
		students := getJSON[[]map[string]any](t, h, "GET", "/api/v1/club/students", coachTok, nil)
		require.NotEmpty(t, students)
		sid := students[0]["id"].(string)
		_ = getJSON[map[string]any](t, h, "DELETE", "/api/v1/club/students/"+sid, coachTok, nil)
	})

	t.Run("leave_and_join", func(t *testing.T) {
		email := fmt.Sprintf("a%d@test.run", time.Now().UnixNano())
		reg := getJSON[map[string]any](t, h, "POST", "/api/v1/auth/register", "", map[string]string{
			"name": "Ath", "email": email, "password": "secret1", "role": "athlete",
		})
		tok := reg["token"].(string)
		_ = getJSON[map[string]any](t, h, "POST", "/api/v1/club/join", tok, map[string]string{"code": "PULSE-7K42"})
		ann := getJSON[[]map[string]any](t, h, "GET", "/api/v1/announces", tok, nil)
		require.NotEmpty(t, ann)
		_ = getJSON[map[string]any](t, h, "POST", "/api/v1/club/leave", tok, nil)
		ann2 := getJSON[[]map[string]any](t, h, "GET", "/api/v1/announces", tok, nil)
		require.Empty(t, ann2)
	})
}

func login(t *testing.T, h http.Handler, email, pass string) string {
	t.Helper()
	res := getJSON[map[string]any](t, h, "POST", "/api/v1/auth/login", "", map[string]string{"email": email, "password": pass})
	return res["token"].(string)
}

func getJSON[T any](t *testing.T, h http.Handler, method, path, token string, body any) T {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Less(t, rr.Code, 400, rr.Body.String())
	var out T
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &out))
	return out
}
