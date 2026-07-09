package auth_usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) tokenResponse(ctx context.Context, user *model.User) (*dto.AuthResponse, error) {
	token, err := u.jwt.Issue(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("jwt.Issue: %w", err)
	}
	view, err := u.toView(ctx, user)
	if err != nil {
		return nil, err
	}
	return dto.NewAuthResponse(token, *view), nil
}

func (u *UseCase) toView(ctx context.Context, user *model.User) (*dto.UserView, error) {
	view := dto.NewUserView(user.ID, user.Name, user.Email, string(user.Role))
	m, err := u.membershipRepo.GetActiveByUser(ctx, user.ID)
	if err == nil {
		view.WithClub(m.ClubID)
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("membershipRepo.GetActiveByUser: %w", err)
	}
	if user.Role == model.RoleCoach {
		club, err := u.clubRepo.GetByCoachID(ctx, user.ID)
		if err == nil {
			return view.WithClub(club.ID), nil
		}
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
		}
		return view.MarkNeedsClub(), nil
	}
	return view, nil
}
