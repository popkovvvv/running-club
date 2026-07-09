package club_usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

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
	planKm, hasPlan := 0.0, false
	pw, err := u.planWeekRepo.GetByClubAndIndex(ctx, club.ID, 0)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("planWeekRepo.GetByClubAndIndex: %w", err)
		}
	} else {
		planKm, hasPlan = pw.TargetKm()
	}
	out := make([]dto.StudentView, 0, len(users))
	for _, usr := range users {
		km, err := u.activityRepo.SumDistByUser(ctx, usr.ID)
		if err != nil {
			return nil, fmt.Errorf("activityRepo.SumDistByUser: %w", err)
		}
		sub, err := u.announceRepo.NextLabelForAthlete(ctx, club.ID, usr.ID)
		if err != nil {
			if !errors.Is(err, model.ErrNotFound) {
				return nil, fmt.Errorf("announceRepo.NextLabelForAthlete: %w", err)
			}
			sub = "Нет записи"
		}
		comp := 0
		if hasPlan {
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
