package workout_usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

const planWeekCount = 4

func (u *UseCase) Plan(ctx context.Context, userID uuid.UUID, week, year, month int) (*dto.PlanResponse, error) {
	if year > 0 && month >= 1 && month <= 12 {
		return u.planMonth(ctx, userID, year, month)
	}
	return u.planWeek(ctx, userID, week)
}

func (u *UseCase) planMonth(ctx context.Context, userID uuid.UUID, year, month int) (*dto.PlanResponse, error) {
	all, err := u.workoutRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUser: %w", err)
	}

	var days, mine []*model.Workout
	for _, w := range all {
		at := workoutDate(w)
		if at.Year() != year || int(at.Month()) != month {
			continue
		}
		switch w.Kind {
		case model.WorkoutPlan:
			days = append(days, w)
		case model.WorkoutOwn:
			mine = append(mine, w)
		}
	}

	weekRange, weekPlan, err := u.currentWeekMeta(ctx, userID)
	if err != nil {
		return nil, err
	}

	weekKm, err := u.activityRepo.SumDistByUserSince(ctx, userID, startOfWeek(time.Now().UTC()))
	if err != nil {
		return nil, fmt.Errorf("activityRepo.SumDistByUserSince: %w", err)
	}

	return dto.NewPlanResponse(0, weekRange, weekPlan, formatKm(weekKm), mapWorkouts(days), mapWorkouts(mine)), nil
}

func (u *UseCase) currentWeekMeta(ctx context.Context, userID uuid.UUID) (weekRange, weekPlan string, err error) {
	_, isoWeek := time.Now().UTC().ISOWeek()
	membership, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return "", "", fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
		}
		return defaultWeekRange(0, time.Now().UTC()), "", nil
	}
	weeks, err := u.planWeekRepo.FindByClub(ctx, membership.ClubID)
	if err != nil {
		return "", "", fmt.Errorf("planWeekRepo.FindByClub: %w", err)
	}
	for _, pw := range weeks {
		if pw.WeekIndex == isoWeek {
			weekRange = pw.RangeLabel
			weekPlan = pw.PlanLabel
			break
		}
	}
	if weekRange == "" {
		weekRange = defaultWeekRange(0, time.Now().UTC())
	}
	if weekPlan == "" {
		planDays, err := u.workoutRepo.FindByUserWeek(ctx, userID, isoWeek, model.WorkoutPlan)
		if err != nil {
			return "", "", fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
		}
		ownDays, err := u.workoutRepo.FindByUserWeek(ctx, userID, isoWeek, model.WorkoutOwn)
		if err != nil {
			return "", "", fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
		}
		weekPlan = planVolumeLabel(planDays, ownDays)
	}
	return weekRange, weekPlan, nil
}

func (u *UseCase) planWeek(ctx context.Context, userID uuid.UUID, week int) (*dto.PlanResponse, error) {
	if week < 0 {
		week = 0
	}
	if week >= planWeekCount {
		week = planWeekCount - 1
	}

	weekRange, weekPlan := defaultWeekRange(week, time.Now().UTC()), ""
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
					if pw.RangeLabel != "" {
						weekRange = pw.RangeLabel
					}
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
	mine, err := u.workoutRepo.FindByUserWeek(ctx, userID, week, model.WorkoutOwn)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
	}
	if weekPlan == "" {
		weekPlan = planVolumeLabel(days, mine)
	}

	weekKm, err := u.activityRepo.SumDistByUserSince(ctx, userID, startOfWeek(time.Now().UTC()))
	if err != nil {
		return nil, fmt.Errorf("activityRepo.SumDistByUserSince: %w", err)
	}

	return dto.NewPlanResponse(week, weekRange, weekPlan, formatKm(weekKm), mapWorkouts(days), mapWorkouts(mine)), nil
}

func workoutDate(w *model.Workout) time.Time {
	if w.ScheduledDate != nil {
		return w.ScheduledDate.UTC()
	}
	return w.CreatedAt.UTC()
}

func defaultWeekRange(weekIndex int, now time.Time) string {
	start := startOfWeek(now).AddDate(0, 0, weekIndex*7)
	end := start.AddDate(0, 0, 6)
	return start.Format("02.01") + " – " + end.Format("02.01")
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func planVolumeLabel(days, mine []*model.Workout) string {
	var total float64
	for _, w := range days {
		total += w.DistKm
	}
	for _, w := range mine {
		total += w.DistKm
	}
	if total <= 0 {
		return ""
	}
	if total == float64(int(total)) {
		return strconv.Itoa(int(total)) + " км"
	}
	return strconv.FormatFloat(total, 'f', 1, 64) + " км"
}

func formatKm(v float64) string {
	if v <= 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}
