package club_usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) resolveClub(ctx context.Context, userID uuid.UUID, role string) (*model.Club, error) {
	if role == string(model.RoleCoach) {
		club, err := u.clubRepo.GetByCoachID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
		}
		return club, nil
	}
	m, err := u.membershipRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrNotMember
		}
		return nil, fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
	}
	club, err := u.clubRepo.GetByID(ctx, m.ClubID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByID: %w", err)
	}
	return club, nil
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
