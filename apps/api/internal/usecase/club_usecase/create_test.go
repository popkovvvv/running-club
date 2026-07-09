//go:build unit

package club_usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/club_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	existing := model.ClubFixture(uuid.New(), "Existing", "PULSE-AAAA", "#ff5c22", coachID)
	wantCode := "PULSE-" + strings.ToUpper(coachID.String()[:4])

	tests := []struct {
		name    string
		req     dto.CreateClubRequest
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			req:  dto.CreateClubRequest{Name: "  Pulse Run  ", AccentHex: "bad"},
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(nil, model.ErrNotFound).Once()
				m.clubRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Club")).Return(nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, mock.Anything).Return(0, nil).Once()
			},
		},
		{
			name: "already_has_club",
			req:  dto.CreateClubRequest{Name: "Pulse Run", AccentHex: "#c8ff34"},
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(existing, nil).Once()
			},
			wantErr: model.ErrConflict,
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
			view, err := uc.Create(context.Background(), coachID, tt.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "Pulse Run", view.Name)
			require.Equal(t, "#ff5c22", view.AccentHex)
			require.Equal(t, wantCode, view.InviteCode)
			require.Equal(t, 0, view.Students)
		})
	}
}
