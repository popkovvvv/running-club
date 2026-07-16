package analytics_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		clubRepo     clubRepo
		userRepo     userRepo
		activityRepo activityRepo
		planWeekRepo planWeekRepo
		workoutRepo  workoutRepo
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	userRepo interface {
		FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error)
	}

	activityRepo interface {
		SumDistByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (float64, error)
		SumDistByClubAthletes(ctx context.Context, clubID uuid.UUID) (float64, error)
	}

	planWeekRepo interface {
		GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error)
	}

	workoutRepo interface {
		SumPlanDistByUserWeek(ctx context.Context, userID uuid.UUID, weekIndex int) (float64, error)
	}
)

func NewUseCase(
	clubRepo clubRepo,
	userRepo userRepo,
	activityRepo activityRepo,
	planWeekRepo planWeekRepo,
	workoutRepo workoutRepo,
) *UseCase {
	return &UseCase{
		clubRepo:     clubRepo,
		userRepo:     userRepo,
		activityRepo: activityRepo,
		planWeekRepo: planWeekRepo,
		workoutRepo:  workoutRepo,
	}
}
