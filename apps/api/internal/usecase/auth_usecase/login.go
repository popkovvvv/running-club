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

func (u *UseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	if !password.Check(user.PasswordHash, req.Password) {
		return nil, model.ErrInvalidCredentials
	}
	return u.tokenResponse(ctx, user)
}
