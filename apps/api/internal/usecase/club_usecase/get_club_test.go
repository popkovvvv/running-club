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

func TestGetClub(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	athleteID := uuid.New()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	mem := model.MembershipFixture(uuid.New(), athleteID, clubID)

	tests := []struct {
		name       string
		userID     uuid.UUID
		role       string
		before     func(m usecaseMocks)
		wantErr    error
		wantInvite string
	}{
		{
			name:   "coach",
			userID: coachID,
			role:   string(model.RoleCoach),
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, clubID).Return(3, nil).Once()
			},
			wantInvite: "PULSE-7K42",
		},
		{
			name:   "athlete",
			userID: athleteID,
			role:   string(model.RoleAthlete),
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(mem, nil).Once()
				m.clubRepo.EXPECT().GetByID(mock.Anything, clubID).Return(club, nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, clubID).Return(3, nil).Once()
			},
		},
		{
			name:   "athlete_not_member",
			userID: athleteID,
			role:   string(model.RoleAthlete),
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
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo)
			view, err := uc.GetClub(context.Background(), tt.userID, tt.role)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, clubID, view.ID)
			require.Equal(t, 3, view.Students)
			require.Equal(t, tt.wantInvite, view.InviteCode)
		})
	}
}
