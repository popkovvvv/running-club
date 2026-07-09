package activity_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) PRs(ctx context.Context, userID uuid.UUID) ([]dto.PRView, error) {
	items, err := u.activityRepo.FindPRs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindPRs: %w", err)
	}
	out := make([]dto.PRView, 0, len(items))
	for _, p := range items {
		out = append(out, dto.NewPRView(p.Distance, p.Time, p.DateLabel))
	}
	return out, nil
}
