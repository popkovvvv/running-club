package club_usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) UpdatePalette(ctx context.Context, coachID uuid.UUID, accent string) (*dto.ClubView, error) {
	accent = strings.TrimSpace(accent)
	if !strings.HasPrefix(accent, "#") || (len(accent) != 7) {
		return nil, fmt.Errorf("%w: invalid accent", model.ErrInvalidInviteCode)
	}
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	if err := u.clubRepo.UpdateAccent(ctx, club.ID, accent); err != nil {
		return nil, fmt.Errorf("clubRepo.UpdateAccent: %w", err)
	}
	club.AccentHex = accent
	n, _ := u.clubRepo.CountActiveStudents(ctx, club.ID)
	return dto.NewClubView(club.ID, club.Name, club.AccentHex, n).WithInviteCode(club.InviteCode), nil
}
