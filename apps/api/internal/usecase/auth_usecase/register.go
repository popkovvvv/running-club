package auth_usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/password"
)

func (u *UseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Password) < 6 {
		return nil, model.ErrWeakPassword
	}
	role := model.Role(req.Role)
	if role != model.RoleAthlete && role != model.RoleCoach {
		return nil, model.ErrInvalidRole
	}
	if _, err := u.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, model.ErrEmailTaken
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password.Hash: %w", err)
	}
	user := model.NewUser(req.Name, req.Email, hash, role)
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("userRepo.Create: %w", err)
	}
	return u.tokenResponse(ctx, user)
}
