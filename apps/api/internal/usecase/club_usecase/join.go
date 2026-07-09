package club_usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Join(ctx context.Context, userID uuid.UUID, code string) (*dto.ClubView, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return nil, model.ErrInvalidInviteCode
	}
	if _, err := u.membershipRepo.GetActiveByUser(ctx, userID); err == nil {
		return nil, model.ErrAlreadyMember
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
	}
	club, err := u.clubRepo.GetByInviteCode(ctx, code)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrInvalidInviteCode
		}
		return nil, fmt.Errorf("clubRepo.GetByInviteCode: %w", err)
	}
	existing, err := u.membershipRepo.GetByUserAndClub(ctx, userID, club.ID)
	if err == nil {
		if err := u.membershipRepo.UpdateStatus(ctx, existing.ID, model.MembershipActive); err != nil {
			return nil, fmt.Errorf("membershipRepo.UpdateStatus: %w", err)
		}
	} else if errors.Is(err, model.ErrNotFound) {
		m := model.NewMembership(userID, club.ID)
		if err := u.membershipRepo.Create(ctx, m); err != nil {
			return nil, fmt.Errorf("membershipRepo.Create: %w", err)
		}
	} else {
		return nil, fmt.Errorf("membershipRepo.GetByUserAndClub: %w", err)
	}
	n, _ := u.clubRepo.CountActiveStudents(ctx, club.ID)
	return dto.NewClubView(club.ID, club.Name, club.AccentHex, n), nil
}
