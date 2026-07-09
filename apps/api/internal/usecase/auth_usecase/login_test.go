//go:build unit

package auth_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
	"github.com/nikpopkov/running-club/api/internal/pkg/password"
	"github.com/nikpopkov/running-club/api/internal/usecase/auth_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	hash, err := password.Hash("password")
	require.NoError(t, err)
	uid := uuid.New()
	user := model.UserFixture(uid, "Никита", "nikita@pulse.run", hash, model.RoleAthlete)

	tests := []struct {
		name    string
		email   string
		pass    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name:  "ok",
			email: "nikita@pulse.run",
			pass:  "password",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "nikita@pulse.run").Return(user, nil).Once()
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
			},
		},
		{
			name:  "wrong_password",
			email: "nikita@pulse.run",
			pass:  "nope",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "nikita@pulse.run").Return(user, nil).Once()
			},
			wantErr: model.ErrInvalidCredentials,
		},
		{
			name:  "unknown",
			email: "x@y.z",
			pass:  "password",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByEmail(mock.Anything, "x@y.z").Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrInvalidCredentials,
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
			res, err := uc.Login(context.Background(), dto.LoginRequest{Email: tt.email, Password: tt.pass})
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, res.Token)
		})
	}
}
