package analytics_usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	signedUp, capacity, err := u.announceRepo.AttendanceStats(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.AttendanceStats: %w", err)
	}
	attendance := attendancePct(signedUp, capacity)
	planKm, hasPlan := 0.0, false
	pw, err := u.planWeekRepo.GetByClubAndIndex(ctx, club.ID, 0)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("planWeekRepo.GetByClubAndIndex: %w", err)
		}
	} else {
		planKm, hasPlan = pw.TargetKm()
	}
	students := make([]dto.StudentView, 0, len(users))
	for _, usr := range users {
		km, err := u.activityRepo.SumDistByUser(ctx, usr.ID)
		if err != nil {
			return nil, fmt.Errorf("activityRepo.SumDistByUser: %w", err)
		}
		comp := 0
		if hasPlan {
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
	return dto.NewAnalyticsResponse(clubKm, attendance, students), nil
}

func attendancePct(signedUp, capacity int) int {
	if signedUp == 0 && capacity == 0 {
		return 0
	}
	denom := capacity
	if denom == 0 {
		denom = signedUp
		if denom < 1 {
			denom = 1
		}
	}
	return int(math.Round(100 * float64(signedUp) / float64(denom)))
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
