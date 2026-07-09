package activity_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		activityRepo activityRepo
	}

	activityRepo interface {
		FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Activity, error)
		FindPRs(ctx context.Context, userID uuid.UUID) ([]*model.PR, error)
		FindRaces(ctx context.Context, userID uuid.UUID) ([]*model.Race, error)
		FindMonthStats(ctx context.Context, userID uuid.UUID) ([]model.MonthStat, error)
	}
)

func NewUseCase(activityRepo activityRepo) *UseCase {
	return &UseCase{activityRepo: activityRepo}
}
