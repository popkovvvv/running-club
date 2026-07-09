package workout_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Plan(ctx context.Context, userID uuid.UUID, week int) (*dto.PlanResponse, error) {
	if week < 0 {
		week = 0
	}
	weekRange, weekPlan := "", ""
	membership, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
		}
	} else {
		weeks, err := u.planWeekRepo.FindByClub(ctx, membership.ClubID)
		if err != nil {
			return nil, fmt.Errorf("planWeekRepo.FindByClub: %w", err)
		}
		if len(weeks) > 0 {
			if week < weeks[0].WeekIndex {
				week = weeks[0].WeekIndex
			}
			last := weeks[len(weeks)-1].WeekIndex
			if week > last {
				week = last
			}
			for _, pw := range weeks {
				if pw.WeekIndex == week {
					weekRange = pw.RangeLabel
					weekPlan = pw.PlanLabel
					break
				}
			}
		}
	}
	days, err := u.workoutRepo.FindByUserWeek(ctx, userID, week, model.WorkoutPlan)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
	}
	mine, err := u.workoutRepo.FindOwnByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindOwnByUser: %w", err)
	}
	return dto.NewPlanResponse(week, weekRange, weekPlan, mapWorkouts(days), mapWorkouts(mine)), nil
}
