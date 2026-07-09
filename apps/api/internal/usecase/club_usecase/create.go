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

func (u *UseCase) Create(ctx context.Context, coachID uuid.UUID, req dto.CreateClubRequest) (*dto.ClubView, error) {
	name := strings.TrimSpace(req.Name)
	accent := strings.TrimSpace(req.AccentHex)
	if name == "" {
		return nil, fmt.Errorf("%w: empty name", model.ErrInvalidInviteCode)
	}
	if !strings.HasPrefix(accent, "#") || len(accent) != 7 {
		accent = "#ff5c22"
	}
	if _, err := u.clubRepo.GetByCoachID(ctx, coachID); err == nil {
		return nil, model.ErrConflict
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	code := "PULSE-" + strings.ToUpper(coachID.String()[:4])
	club := model.NewClub(name, code, accent, coachID)
	if err := u.clubRepo.Create(ctx, club); err != nil {
		return nil, fmt.Errorf("clubRepo.Create: %w", err)
	}
	n, err := u.clubRepo.CountActiveStudents(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.CountActiveStudents: %w", err)
	}
	return dto.NewClubView(club.ID, club.Name, club.AccentHex, n).WithInviteCode(club.InviteCode), nil
}
