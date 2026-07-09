package schedule_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) List(ctx context.Context, userID uuid.UUID, role string) ([]dto.AnnounceView, error) {
	clubID, err := u.clubIDFor(ctx, userID, role)
	if err != nil {
		if errors.Is(err, model.ErrNotMember) || errors.Is(err, model.ErrNotFound) {
			return []dto.AnnounceView{}, nil
		}
		return nil, err
	}
	items, err := u.announceRepo.FindByClub(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.FindByClub: %w", err)
	}
	out := make([]dto.AnnounceView, 0, len(items))
	for _, a := range items {
		signed := false
		if role == string(model.RoleAthlete) {
			signed, err = u.announceRepo.HasSignup(ctx, a.ID, userID)
			if err != nil {
				return nil, fmt.Errorf("announceRepo.HasSignup: %w", err)
			}
		}
		out = append(out, toAnnounceView(a, signed))
	}
	return out, nil
}
