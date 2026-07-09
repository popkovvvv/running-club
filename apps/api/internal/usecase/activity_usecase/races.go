package activity_usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) Races(ctx context.Context, userID uuid.UUID) ([]dto.RaceView, error) {
	items, err := u.activityRepo.FindRaces(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindRaces: %w", err)
	}
	out := make([]dto.RaceView, 0, len(items))
	for _, r := range items {
		if r.Finished {
			continue
		}
		out = append(out, dto.NewRaceView(strconv.Itoa(r.DaysLeft), r.Name, r.DateLabel, r.Dist, r.Goal))
	}
	return out, nil
}
