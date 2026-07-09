package auth_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) Me(ctx context.Context, userID uuid.UUID) (*dto.UserView, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetByID: %w", err)
	}
	view, err := u.toView(ctx, user)
	if err != nil {
		return nil, err
	}
	return view, nil
}
