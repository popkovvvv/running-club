package club_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) GetClub(ctx context.Context, userID uuid.UUID, role string) (*dto.ClubView, error) {
	club, err := u.resolveClub(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	n, err := u.clubRepo.CountActiveStudents(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.CountActiveStudents: %w", err)
	}
	view := dto.NewClubView(club.ID, club.Name, club.AccentHex, n)
	if role == string(model.RoleCoach) {
		view.WithInviteCode(club.InviteCode)
	}
	return view, nil
}
