package analytics_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		clubRepo     clubRepo
		userRepo     userRepo
		activityRepo activityRepo
		announceRepo announceRepo
		planWeekRepo planWeekRepo
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	userRepo interface {
		FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error)
	}

	activityRepo interface {
		SumDistByUser(ctx context.Context, userID uuid.UUID) (float64, error)
		SumDistByClubAthletes(ctx context.Context, clubID uuid.UUID) (float64, error)
	}

	announceRepo interface {
		AttendanceStats(ctx context.Context, clubID uuid.UUID) (signedUp int, capacity int, err error)
	}

	planWeekRepo interface {
		GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error)
	}
)

func NewUseCase(
	clubRepo clubRepo,
	userRepo userRepo,
	activityRepo activityRepo,
	announceRepo announceRepo,
	planWeekRepo planWeekRepo,
) *UseCase {
	return &UseCase{
		clubRepo:     clubRepo,
		userRepo:     userRepo,
		activityRepo: activityRepo,
		announceRepo: announceRepo,
		planWeekRepo: planWeekRepo,
	}
}
