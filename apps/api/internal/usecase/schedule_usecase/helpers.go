package schedule_usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

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

func (u *UseCase) clubFor(ctx context.Context, userID uuid.UUID, role string) (*model.Club, error) {
	if role == string(model.RoleCoach) {
		club, err := u.clubRepo.GetByCoachID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
		}
		return club, nil
	}
	m, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	club, err := u.clubRepo.GetByID(ctx, m.ClubID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByID: %w", err)
	}
	return club, nil
}

func (u *UseCase) toAnnounceView(ctx context.Context, a *model.Announce, signed bool) (dto.AnnounceView, error) {
	cta := "Записаться"
	if signed {
		cta = "Вы записаны"
	}
	v := dto.NewAnnounceView(a.ID, a.Place, a.DayLabel, a.Time, a.GroupName, a.Note, a.GoingCount, signed, cta)
	if a.StartsOn != nil {
		v = v.WithStartsOn(a.StartsOn.Format("2006-01-02"))
	}
	athletes, err := u.announceRepo.FindGoingAthletes(ctx, a.ID)
	if err != nil {
		return dto.AnnounceView{}, fmt.Errorf("announceRepo.FindGoingAthletes: %w", err)
	}
	attendees := make([]dto.GoingPersonView, 0, len(athletes))
	for _, usr := range athletes {
		attendees = append(attendees, dto.NewGoingPersonView(usr.ID, initials(usr.Name), usr.Name))
	}
	return v.WithAttendees(attendees), nil
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
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
