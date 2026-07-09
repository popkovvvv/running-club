package schedule_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Publish(ctx context.Context, coachID uuid.UUID, req dto.CreateAnnounceRequest) (*dto.AnnounceView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	a := model.NewAnnounce(club.ID, req.Place, req.Day, req.Time, req.Group, req.Note)
	if err := u.announceRepo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("announceRepo.Create: %w", err)
	}
	v := toAnnounceView(a, false)
	return &v, nil
}
