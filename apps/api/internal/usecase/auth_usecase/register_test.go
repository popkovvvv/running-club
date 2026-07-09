//go:build unit

package auth_usecase_test

import (
	"context"
	"testing"

	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
	"github.com/nikpopkov/running-club/api/internal/usecase/auth_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		req        dto.RegisterRequest
		before     func(m usecaseMocks)
		wantErr    error
		wantRole   string
		wantNeeds  *bool
		wantInClub *bool
	}{
		{
			name: "ok_athlete",
			req:  dto.RegisterRequest{Name: "Никита", Email: "a@b.c", Password: "secret1", Role: "athlete"},
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "a@b.c").Return(nil, model.ErrNotFound).Once()
				m.userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()
			},
			wantRole: "athlete",
		},
		{
			name: "ok_coach_no_auto_club",
			req:  dto.RegisterRequest{Name: "Coach", Email: "coach@b.c", Password: "secret1", Role: "coach"},
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "coach@b.c").Return(nil, model.ErrNotFound).Once()
				m.userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()
			},
			wantRole:   "coach",
			wantNeeds:  boolPtr(true),
			wantInClub: boolPtr(false),
		},
		{
			name:    "weak_password",
			req:     dto.RegisterRequest{Name: "X", Email: "a@b.c", Password: "123", Role: "athlete"},
			wantErr: model.ErrWeakPassword,
		},
		{
			name:    "invalid_role",
			req:     dto.RegisterRequest{Name: "X", Email: "a@b.c", Password: "secret1", Role: "admin"},
			wantErr: model.ErrInvalidRole,
		},
		{
			name: "email_taken",
			req:  dto.RegisterRequest{Name: "X", Email: "taken@b.c", Password: "secret1", Role: "athlete"},
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "taken@b.c").
					Return(model.NewUser("X", "taken@b.c", "", model.RoleAthlete), nil).Once()
			},
			wantErr: model.ErrEmailTaken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := auth_usecase.NewUseCase(m.userRepo, m.membershipRepo, m.clubRepo, authjwt.NewManager("test"))
			res, err := uc.Register(context.Background(), tt.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, res.Token)
			require.Equal(t, tt.wantRole, res.User.Role)
			if tt.wantNeeds != nil {
				require.Equal(t, *tt.wantNeeds, res.User.NeedsClub)
			}
			if tt.wantInClub != nil {
				require.Equal(t, *tt.wantInClub, res.User.InClub)
			}
		})
	}
}

func boolPtr(v bool) *bool {
	return &v
}
