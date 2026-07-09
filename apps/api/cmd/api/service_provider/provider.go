package service_provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/activity_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/announce_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/membership_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/plan_week_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/workout_repo"
	"github.com/nikpopkov/running-club/api/internal/app/http/activity"
	"github.com/nikpopkov/running-club/api/internal/app/http/analytics"
	"github.com/nikpopkov/running-club/api/internal/app/http/auth"
	"github.com/nikpopkov/running-club/api/internal/app/http/club"
	"github.com/nikpopkov/running-club/api/internal/app/http/router"
	"github.com/nikpopkov/running-club/api/internal/app/http/schedule"
	"github.com/nikpopkov/running-club/api/internal/app/http/workout"
	"github.com/nikpopkov/running-club/api/internal/config"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/analytics_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/auth_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/club_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase"
)

type ServiceProvider struct {
	cfg  config.Config
	pool *pgxpool.Pool
	jwt  *authjwt.Manager
}

func New(cfg config.Config) *ServiceProvider {
	return &ServiceProvider{cfg: cfg}
}

func (s *ServiceProvider) Boot(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, s.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("pool.Ping: %w", err)
	}
	s.pool = pool
	s.jwt = authjwt.NewManager(s.cfg.JWTSecret)
	return nil
}

func (s *ServiceProvider) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *ServiceProvider) Pool() *pgxpool.Pool { return s.pool }
func (s *ServiceProvider) JWT() *authjwt.Manager { return s.jwt }
func (s *ServiceProvider) Config() config.Config { return s.cfg }

func (s *ServiceProvider) Handler() http.Handler {
	userRepo := user_repo.NewRepo(s.pool)
	clubRepo := club_repo.NewRepo(s.pool)
	membershipRepo := membership_repo.NewRepo(s.pool)
	announceRepo := announce_repo.NewRepo(s.pool)
	workoutRepo := workout_repo.NewRepo(s.pool)
	activityRepo := activity_repo.NewRepo(s.pool)
	planWeekRepo := plan_week_repo.NewRepo(s.pool)

	authUC := auth_usecase.NewUseCase(userRepo, membershipRepo, clubRepo, s.jwt)
	clubUC := club_usecase.NewUseCase(clubRepo, membershipRepo, userRepo, activityRepo, announceRepo, planWeekRepo)
	scheduleUC := schedule_usecase.NewUseCase(announceRepo, clubRepo, membershipRepo)
	workoutUC := workout_usecase.NewUseCase(workoutRepo, planWeekRepo, membershipRepo)
	activityUC := activity_usecase.NewUseCase(activityRepo)
	analyticsUC := analytics_usecase.NewUseCase(clubRepo, userRepo, activityRepo, announceRepo, planWeekRepo)

	return router.New(router.Handlers{
		Auth:      auth.NewHandler(authUC),
		Club:      club.NewHandler(clubUC),
		Schedule:  schedule.NewHandler(scheduleUC),
		Workout:   workout.NewHandler(workoutUC),
		Activity:  activity.NewHandler(activityUC),
		Analytics: analytics.NewHandler(analyticsUC),
		JWT:       s.jwt,
	})
}
