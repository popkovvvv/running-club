package analytics_usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
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
	students := make([]dto.StudentView, 0, len(users))
	for i, usr := range users {
		comp := 80
		if i%3 == 0 {
			comp = 90
		}
		students = append(students, dto.NewStudentView(usr.ID, initials(usr.Name), usr.Name, "Прогресс недели", "24.6", comp))
	}
	return dto.NewAnalyticsResponse(186.1, 86, students), nil
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
