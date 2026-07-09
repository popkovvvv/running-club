package club_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) RemoveStudent(ctx context.Context, coachID, studentID uuid.UUID) error {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	m, err := u.membershipRepo.GetByUserAndClub(ctx, studentID, club.ID)
	if err != nil {
		return fmt.Errorf("membershipRepo.GetByUserAndClub: %w", err)
	}
	if err := u.membershipRepo.UpdateStatus(ctx, m.ID, model.MembershipRemoved); err != nil {
		return fmt.Errorf("membershipRepo.UpdateStatus: %w", err)
	}
	return nil
}
