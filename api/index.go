package handler

import (
	"net/http"

	"github.com/nikpopkov/running-club/api/vercelentry"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	vercelentry.Handler(w, r)
}
