package workout_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Plan(ctx context.Context, userID uuid.UUID, week int) (*dto.PlanResponse, error) {
	if week < 0 {
		week = 0
	}
	if week >= len(weekMeta) {
		week = len(weekMeta) - 1
	}
	days, err := u.workoutRepo.FindByUserWeek(ctx, userID, week, model.WorkoutPlan)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
	}
	mine, err := u.workoutRepo.FindOwnByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindOwnByUser: %w", err)
	}
	return dto.NewPlanResponse(week, weekMeta[week].Range, weekMeta[week].Plan, mapWorkouts(days), mapWorkouts(mine)), nil
}
