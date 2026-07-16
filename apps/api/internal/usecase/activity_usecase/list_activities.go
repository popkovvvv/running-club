package activity_usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) ListActivities(ctx context.Context, userID uuid.UUID) ([]dto.ActivityView, error) {
	items, err := u.activityRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindByUser: %w", err)
	}
	out := make([]dto.ActivityView, 0, len(items))
	for _, a := range items {
		out = append(out, dto.NewActivityView(
			a.ID, a.Title, a.WhenLabel, formatKm(a.DistKm), a.Duration, a.Pace, strconv.Itoa(a.HR),
			a.Kudos, a.Comments, a.RouteSVG, a.StartX, a.StartY, a.EndX, a.EndY,
			a.Source, a.SportType, a.ElevationGain, a.Visibility,
		))
	}
	return out, nil
}

func formatKm(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
