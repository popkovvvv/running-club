package schedule_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Publish(ctx context.Context, coachID uuid.UUID, req dto.CreateAnnounceRequest) (*dto.AnnounceView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	var startsOn *time.Time
	if req.StartsOn != "" {
		t, err := time.Parse("2006-01-02", req.StartsOn)
		if err != nil {
			return nil, fmt.Errorf("parse startsOn: %w", err)
		}
		startsOn = &t
	}
	a := model.NewAnnounce(club.ID, req.Place, req.Day, req.Time, req.Group, req.Note, startsOn)
	if err := u.announceRepo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("announceRepo.Create: %w", err)
	}
	v, err := u.toAnnounceView(ctx, a, false)
	if err != nil {
		return nil, fmt.Errorf("toAnnounceView: %w", err)
	}
	return &v, nil
}
