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
	views := make([]dto.MonthStatView, 0, len(months))
	for _, m := range months {
		views = append(views, dto.NewMonthStatView(m.Month, m.Km, m.Tr, m.Pace, m.Diff))
	}
	return dto.NewProgressResponse(612, 78, 3, views), nil
}
