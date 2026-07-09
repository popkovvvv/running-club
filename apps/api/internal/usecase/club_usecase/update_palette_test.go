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

func TestUpdatePalette(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)

	tests := []struct {
		name    string
		accent  string
		before  func(m usecaseMocks)
		wantErr bool
	}{
		{
			name:   "ok",
			accent: "#c8ff34",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.clubRepo.EXPECT().UpdateAccent(mock.Anything, clubID, "#c8ff34").Return(nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, clubID).Return(2, nil).Once()
			},
		},
		{
			name:    "invalid_accent",
			accent:  "red",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo, m.activityRepo, m.announceRepo)
			view, err := uc.UpdatePalette(context.Background(), coachID, tt.accent)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "#c8ff34", view.AccentHex)
			require.Equal(t, "PULSE-7K42", view.InviteCode)
		})
	}
}
