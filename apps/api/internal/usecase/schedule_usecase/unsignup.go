package schedule_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Unsignup(ctx context.Context, athleteID, announceID uuid.UUID) (*dto.AnnounceView, error) {
	a, err := u.announceRepo.GetByID(ctx, announceID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.GetByID: %w", err)
	}
	if err := u.announceRepo.DeleteSignup(ctx, announceID, athleteID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrNotSignedUp
		}
		return nil, fmt.Errorf("announceRepo.DeleteSignup: %w", err)
	}
	if err := u.announceRepo.IncGoing(ctx, announceID, -1); err != nil {
		return nil, fmt.Errorf("announceRepo.IncGoing: %w", err)
	}
	if err := u.workoutRepo.DeleteByUserAndAnnounce(ctx, athleteID, announceID); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("workoutRepo.DeleteByUserAndAnnounce: %w", err)
		}
	}
	a, err = u.announceRepo.GetByID(ctx, announceID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.GetByID: %w", err)
	}
	v, err := u.toAnnounceView(ctx, a, false)
	if err != nil {
		return nil, fmt.Errorf("toAnnounceView: %w", err)
	}
	return &v, nil
}
