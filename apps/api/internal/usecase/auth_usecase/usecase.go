package auth_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
)

type (
	UseCase struct {
		userRepo       userRepo
		membershipRepo membershipRepo
		clubRepo       clubCreator
		jwt            *authjwt.Manager
	}

	userRepo interface {
		Create(ctx context.Context, u *model.User) error
		GetByEmail(ctx context.Context, email string) (*model.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	}

	membershipRepo interface {
		GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
	}

	clubCreator interface {
		Create(ctx context.Context, c *model.Club) error
	}
)

func NewUseCase(userRepo userRepo, membershipRepo membershipRepo, clubRepo clubCreator, jwt *authjwt.Manager) *UseCase {
	return &UseCase{userRepo: userRepo, membershipRepo: membershipRepo, clubRepo: clubRepo, jwt: jwt}
}
