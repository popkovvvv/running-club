package workout_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		workoutRepo    workoutRepo
		planWeekRepo   planWeekRepo
		membershipRepo membershipRepo
		clubRepo       clubRepo
		activityRepo   activityRepo
	}

	workoutRepo interface {
		Create(ctx context.Context, w *model.Workout) error
		GetByID(ctx context.Context, id uuid.UUID) (*model.Workout, error)
		FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error)
		FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error)
		Delete(ctx context.Context, id uuid.UUID) error
		Update(ctx context.Context, w *model.Workout) error
	}

	activityRepo interface {
		Create(ctx context.Context, a *model.Activity) error
		GetByID(ctx context.Context, id uuid.UUID) (*model.Activity, error)
		SumDistByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (float64, error)
	}

	planWeekRepo interface {
		FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.PlanWeek, error)
	}

	membershipRepo interface {
		GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
		GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error)
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}
)

func NewUseCase(
	workoutRepo workoutRepo,
	planWeekRepo planWeekRepo,
	membershipRepo membershipRepo,
	clubRepo clubRepo,
	activityRepo activityRepo,
) *UseCase {
	return &UseCase{
		workoutRepo:    workoutRepo,
		planWeekRepo:   planWeekRepo,
		membershipRepo: membershipRepo,
		clubRepo:       clubRepo,
		activityRepo:   activityRepo,
	}
}
