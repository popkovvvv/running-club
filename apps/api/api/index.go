package handler

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/nikpopkov/running-club/api/cmd/api/service_provider"
	"github.com/nikpopkov/running-club/api/internal/config"
)

var (
	once sync.Once
	h    http.Handler
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		cfg := config.Load()
		sp := service_provider.New(cfg)
		if err := sp.Boot(context.Background()); err != nil {
			panic(err)
		}
		h = sp.Handler()
		_ = os.Setenv("SEEDED", "1")
	})
	h.ServeHTTP(w, r)
}
