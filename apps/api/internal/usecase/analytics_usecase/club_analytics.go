package analytics_usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) ClubAnalytics(ctx context.Context, coachID uuid.UUID) (*dto.AnalyticsResponse, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	users, err := u.userRepo.FindAthletesByClub(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.FindAthletesByClub: %w", err)
	}
	clubKm, err := u.activityRepo.SumDistByClubAthletes(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.SumDistByClubAthletes: %w", err)
	}
	weekStart := startOfWeek(time.Now().UTC())
	students := make([]dto.StudentView, 0, len(users))
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
		comp := 0
		if planKm > 0 {
			comp = int(math.Min(100, math.Round(100*km/planKm)))
		}
		students = append(students, dto.NewStudentView(
			usr.ID,
			initials(usr.Name),
			usr.Name,
			"Прогресс недели",
			strconv.FormatFloat(km, 'f', 1, 64),
			comp,
		))
	}
	return dto.NewAnalyticsResponse(clubKm, students), nil
}

func initials(name string) string {
	parts := strings.Fields(name)
	var b strings.Builder
	for i, p := range parts {
		if i >= 2 {
			break
		}
		r, _ := utf8.DecodeRuneInString(p)
		b.WriteRune(r)
	}
	return strings.ToUpper(b.String())
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
