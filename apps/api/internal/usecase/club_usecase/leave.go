package club_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Leave(ctx context.Context, userID uuid.UUID) error {
	m, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.ErrNotMember
		}
		return fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
	}
	if err := u.membershipRepo.UpdateStatus(ctx, m.ID, model.MembershipLeft); err != nil {
		return fmt.Errorf("membershipRepo.UpdateStatus: %w", err)
	}
	return nil
}
