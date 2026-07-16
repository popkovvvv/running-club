package plan_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		planWeekRepo   planWeekRepo
		workoutRepo    workoutRepo
		clubRepo       clubRepo
		userRepo       userRepo
		membershipRepo membershipRepo
	}

	planWeekRepo interface {
		FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.PlanWeek, error)
		GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error)
		Upsert(ctx context.Context, w *model.PlanWeek) error
	}

	workoutRepo interface {
		FindClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int) ([]*model.Workout, error)
		ReplaceClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int, workouts []*model.Workout) error
		DeleteClubAssignedPlans(ctx context.Context, userID uuid.UUID, weekIndex int) error
		Create(ctx context.Context, w *model.Workout) error
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	userRepo interface {
		FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	}

	membershipRepo interface {
		GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error)
	}
)

func NewUseCase(
	planWeekRepo planWeekRepo,
	workoutRepo workoutRepo,
	clubRepo clubRepo,
	userRepo userRepo,
	membershipRepo membershipRepo,
) *UseCase {
	return &UseCase{
		planWeekRepo:   planWeekRepo,
		workoutRepo:    workoutRepo,
		clubRepo:       clubRepo,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
	}
}
