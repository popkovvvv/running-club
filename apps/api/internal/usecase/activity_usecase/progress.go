package activity_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) Progress(ctx context.Context, userID uuid.UUID) (*dto.ProgressResponse, error) {
	months, err := u.activityRepo.FindMonthStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindMonthStats: %w", err)
	}
	var yearKm float64
	var yearTr int
	views := make([]dto.MonthStatView, 0, len(months))
	for _, m := range months {
		yearKm += m.Km
		yearTr += m.Tr
		views = append(views, dto.NewMonthStatView(m.Month, m.Km, m.Tr, m.Pace, m.Diff))
	}
	races, err := u.activityRepo.FindRaces(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindRaces: %w", err)
	}
	return dto.NewProgressResponse(yearKm, yearTr, len(races), views), nil
}
