package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/app/http/response"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
)

type ctxKey string

const (
	UserIDKey ctxKey = "userID"
	RoleKey   ctxKey = "role"
)

func Auth(jwtMgr *authjwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				response.Error(w, model.ErrUnauthorized)
				return
			}
			claims, err := jwtMgr.Parse(strings.TrimPrefix(h, "Bearer "))
			if err != nil {
				response.Error(w, model.ErrUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) uuid.UUID {
	v, _ := ctx.Value(UserIDKey).(uuid.UUID)
	return v
}

func Role(ctx context.Context) string {
	v, _ := ctx.Value(RoleKey).(string)
	return v
}
