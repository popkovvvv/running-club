package schedule_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) clubIDFor(ctx context.Context, userID uuid.UUID, role string) (uuid.UUID, error) {
	if role == string(model.RoleCoach) {
		club, err := u.clubRepo.GetByCoachID(ctx, userID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
		}
		return club.ID, nil
	}
	m, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}
	return m.ClubID, nil
}

func toAnnounceView(a *model.Announce, signed bool) dto.AnnounceView {
	cta := "Записаться"
	if signed {
		cta = "Вы записаны"
	}
	return dto.NewAnnounceView(a.ID, a.Place, a.DayLabel, a.Time, a.GroupName, a.Note, a.GoingCount, signed, cta)
}
