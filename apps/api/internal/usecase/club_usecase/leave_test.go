//go:build unit

package club_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/club_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLeave(t *testing.T) {
	t.Parallel()
	athleteID := uuid.New()
	clubID := uuid.New()
	memID := uuid.New()
	mem := model.MembershipFixture(memID, athleteID, clubID)

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(mem, nil).Once()
				m.membershipRepo.EXPECT().UpdateStatus(mock.Anything, memID, model.MembershipLeft).Return(nil).Once()
			},
		},
		{
			name: "not_member",
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrNotMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo, m.activityRepo, m.announceRepo, m.planWeekRepo)
			err := uc.Leave(context.Background(), athleteID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
