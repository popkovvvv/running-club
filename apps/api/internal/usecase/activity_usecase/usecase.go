package activity_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		activityRepo       activityRepo
		activityStreamRepo activityStreamRepo
		workoutRepo        workoutRepo
		clubRepo           clubRepo
		membershipRepo     membershipRepo
	}

	activityRepo interface {
		FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Activity, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Activity, error)
		FindPRs(ctx context.Context, userID uuid.UUID) ([]*model.PR, error)
		FindRaces(ctx context.Context, userID uuid.UUID) ([]*model.Race, error)
	}

	activityStreamRepo interface {
		FindByActivityID(ctx context.Context, activityID uuid.UUID) ([]*model.ActivityStream, error)
	}

	workoutRepo interface {
		FindByCompletedActivity(ctx context.Context, activityID uuid.UUID) (*model.Workout, error)
		FindCompletedWithoutActivity(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error)
		FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error)
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	membershipRepo interface {
		GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error)
	}
)

func NewUseCase(
	activityRepo activityRepo,
	activityStreamRepo activityStreamRepo,
	workoutRepo workoutRepo,
	clubRepo clubRepo,
	membershipRepo membershipRepo,
) *UseCase {
	return &UseCase{
		activityRepo:       activityRepo,
		activityStreamRepo: activityStreamRepo,
		workoutRepo:        workoutRepo,
		clubRepo:           clubRepo,
		membershipRepo:     membershipRepo,
	}
}
