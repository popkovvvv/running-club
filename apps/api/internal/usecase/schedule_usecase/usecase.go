package schedule_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		announceRepo   announceRepo
		clubRepo       clubResolver
		membershipRepo membershipRepo
	}

	announceRepo interface {
		Create(ctx context.Context, a *model.Announce) error
		FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.Announce, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Announce, error)
		IncGoing(ctx context.Context, id uuid.UUID, delta int) error
		CreateSignup(ctx context.Context, s *model.AnnounceSignup) error
		DeleteSignup(ctx context.Context, announceID, athleteID uuid.UUID) error
		HasSignup(ctx context.Context, announceID, athleteID uuid.UUID) (bool, error)
	}

	clubResolver interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	membershipRepo interface {
		GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
	}
)

func NewUseCase(announceRepo announceRepo, clubRepo clubResolver, membershipRepo membershipRepo) *UseCase {
	return &UseCase{announceRepo: announceRepo, clubRepo: clubRepo, membershipRepo: membershipRepo}
}
