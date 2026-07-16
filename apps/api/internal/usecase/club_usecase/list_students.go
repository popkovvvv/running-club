package club_usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) ListStudents(ctx context.Context, coachID uuid.UUID) ([]dto.StudentView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	users, err := u.userRepo.FindAthletesByClub(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.FindAthletesByClub: %w", err)
	}
	weekStart := startOfWeek(time.Now().UTC())
	out := make([]dto.StudentView, 0, len(users))
	for _, usr := range users {
		km, err := u.activityRepo.SumDistByUserSince(ctx, usr.ID, weekStart)
		if err != nil {
			return nil, fmt.Errorf("activityRepo.SumDistByUserSince: %w", err)
		}
		planKm, err := u.workoutRepo.SumPlanDistByUserWeek(ctx, usr.ID, 0)
		if err != nil {
			return nil, fmt.Errorf("workoutRepo.SumPlanDistByUserWeek: %w", err)
		}
		if planKm == 0 {
			if pw, err := u.planWeekRepo.GetByClubAndIndex(ctx, club.ID, 0); err == nil {
				planKm, _ = pw.TargetKm()
			} else if !errors.Is(err, model.ErrNotFound) {
				return nil, fmt.Errorf("planWeekRepo.GetByClubAndIndex: %w", err)
			}
		}
		sub, err := u.announceRepo.NextLabelForAthlete(ctx, club.ID, usr.ID)
		if err != nil {
			if !errors.Is(err, model.ErrNotFound) {
				return nil, fmt.Errorf("announceRepo.NextLabelForAthlete: %w", err)
			}
			sub = "Нет записи"
		}
		comp := 0
		if planKm > 0 {
			comp = int(math.Min(100, math.Round(100*km/planKm)))
		}
		out = append(out, dto.NewStudentView(
			usr.ID,
			initials(usr.Name),
			usr.Name,
			sub,
			strconv.FormatFloat(km, 'f', 1, 64),
			comp,
		))
	}
	return out, nil
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
