package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nikpopkov/running-club/api/internal/app/http/activity"
	"github.com/nikpopkov/running-club/api/internal/app/http/analytics"
	"github.com/nikpopkov/running-club/api/internal/app/http/auth"
	"github.com/nikpopkov/running-club/api/internal/app/http/club"
	"github.com/nikpopkov/running-club/api/internal/app/http/middleware"
	"github.com/nikpopkov/running-club/api/internal/app/http/schedule"
	"github.com/nikpopkov/running-club/api/internal/app/http/workout"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
)

type Handlers struct {
	Auth      *auth.Handler
	Club      *club.Handler
	Schedule  *schedule.Handler
	Workout   *workout.Handler
	Activity  *activity.Handler
	Analytics *analytics.Handler
	JWT       *authjwt.Manager
}

func New(h Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/register", h.Auth.Register)
		api.Post("/auth/login", h.Auth.Login)

		api.Group(func(priv chi.Router) {
			priv.Use(middleware.Auth(h.JWT))
			priv.Get("/auth/me", h.Auth.Me)
			priv.Post("/auth/logout", h.Auth.Logout)

			priv.Post("/clubs", h.Club.Create)
			priv.Get("/club", h.Club.Get)
			priv.Post("/club/join", h.Club.Join)
			priv.Post("/club/leave", h.Club.Leave)
			priv.Patch("/club/palette", h.Club.Palette)
			priv.Get("/club/students", h.Club.Students)
			priv.Delete("/club/students/{id}", h.Club.RemoveStudent)

			priv.Get("/announces", h.Schedule.List)
			priv.Post("/announces", h.Schedule.Publish)
			priv.Post("/announces/{id}/signup", h.Schedule.Signup)
			priv.Delete("/announces/{id}/signup", h.Schedule.Unsignup)
			priv.Get("/schedule/calendar", h.Schedule.Calendar)

			priv.Get("/plan", h.Workout.Plan)
			priv.Post("/workouts", h.Workout.Create)
			priv.Delete("/workouts/{id}", h.Workout.Delete)

			priv.Get("/activities", h.Activity.List)
			priv.Get("/progress", h.Activity.Progress)
			priv.Get("/prs", h.Activity.PRs)
			priv.Get("/races", h.Activity.Races)
			priv.Get("/analytics", h.Analytics.Get)
		})
	})

	return r
}
