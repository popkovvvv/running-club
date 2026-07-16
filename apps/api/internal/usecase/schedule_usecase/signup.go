package schedule_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Signup(ctx context.Context, athleteID, announceID uuid.UUID) (*dto.AnnounceView, error) {
	a, err := u.announceRepo.GetByID(ctx, announceID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.GetByID: %w", err)
	}
	ok, err := u.announceRepo.HasSignup(ctx, announceID, athleteID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.HasSignup: %w", err)
	}
	if ok {
		return nil, model.ErrAlreadySignedUp
	}
	s := model.NewAnnounceSignup(announceID, athleteID)
	if err := u.announceRepo.CreateSignup(ctx, s); err != nil {
		return nil, fmt.Errorf("announceRepo.CreateSignup: %w", err)
	}
	if err := u.announceRepo.IncGoing(ctx, announceID, 1); err != nil {
		return nil, fmt.Errorf("announceRepo.IncGoing: %w", err)
	}
	w := workoutFromAnnounce(athleteID, a)
	if err := u.workoutRepo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("workoutRepo.Create: %w", err)
	}
	a, err = u.announceRepo.GetByID(ctx, announceID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.GetByID: %w", err)
	}
	v, err := u.toAnnounceView(ctx, a, true)
	if err != nil {
		return nil, fmt.Errorf("toAnnounceView: %w", err)
	}
	return &v, nil
}
