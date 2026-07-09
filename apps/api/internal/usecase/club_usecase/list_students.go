package club_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
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
	out := make([]dto.StudentView, 0, len(users))
	for _, usr := range users {
		out = append(out, dto.NewStudentView(usr.ID, initials(usr.Name), usr.Name, "В клубе", "0", 80))
	}
	return out, nil
}
