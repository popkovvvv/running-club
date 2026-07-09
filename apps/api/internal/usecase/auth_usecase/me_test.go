//go:build unit

package auth_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/authjwt"
	"github.com/nikpopkov/running-club/api/internal/usecase/auth_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMe(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	clubID := uuid.New()
	user := model.UserFixture(uid, "Никита", "nikita@pulse.run", "hash", model.RoleAthlete)
	mem := model.MembershipFixture(uuid.New(), uid, clubID)

	tests := []struct {
		name     string
		before   func(m usecaseMocks)
		wantErr  error
		wantClub *uuid.UUID
	}{
		{
			name: "ok_with_club",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByID(mock.Anything, uid).Return(user, nil).Once()
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(mem, nil).Once()
			},
			wantClub: &clubID,
		},
		{
			name: "ok_without_club",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByID(mock.Anything, uid).Return(user, nil).Once()
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
			},
		},
		{
			name: "not_found",
			before: func(m usecaseMocks) {
				m.userRepo.EXPECT().GetByID(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrNotFound,
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
			view, err := uc.Me(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, uid, view.ID)
			require.Equal(t, "Никита", view.Name)
			if tt.wantClub != nil {
				require.NotNil(t, view.ClubID)
				require.Equal(t, *tt.wantClub, *view.ClubID)
			}
		})
	}
}
