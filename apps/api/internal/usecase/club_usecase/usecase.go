package club_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		clubRepo       clubRepo
		membershipRepo membershipRepo
		userRepo       userRepo
	}

	clubRepo interface {
		Create(ctx context.Context, c *model.Club) error
		GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error)
		GetByInviteCode(ctx context.Context, code string) (*model.Club, error)
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
		UpdateAccent(ctx context.Context, id uuid.UUID, accent string) error
		CountActiveStudents(ctx context.Context, clubID uuid.UUID) (int, error)
	}

	membershipRepo interface {
		Create(ctx context.Context, m *model.Membership) error
		GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
		GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error)
		UpdateStatus(ctx context.Context, id uuid.UUID, status model.MembershipStatus) error
	}

	userRepo interface {
		GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
		FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error)
	}
)

func NewUseCase(clubRepo clubRepo, membershipRepo membershipRepo, userRepo userRepo) *UseCase {
	return &UseCase{clubRepo: clubRepo, membershipRepo: membershipRepo, userRepo: userRepo}
}
