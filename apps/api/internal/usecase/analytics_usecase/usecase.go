package analytics_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		clubRepo clubRepo
		userRepo userRepo
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	userRepo interface {
		FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error)
	}
)

func NewUseCase(clubRepo clubRepo, userRepo userRepo) *UseCase {
	return &UseCase{clubRepo: clubRepo, userRepo: userRepo}
}
