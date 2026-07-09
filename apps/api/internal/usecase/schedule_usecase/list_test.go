//go:build unit

package schedule_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	athleteID := uuid.New()
	clubID := uuid.New()
	annID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	announce := model.AnnounceFixture(annID, clubID, "Зина", "Вт", "19:50", "Основная")
	mem := model.MembershipFixture(uuid.New(), athleteID, clubID)

	tests := []struct {
		name       string
		userID     uuid.UUID
		role       string
		before     func(m usecaseMocks)
		wantLen    int
		wantSigned bool
	}{
		{
			name:   "coach",
			userID: coachID,
			role:   string(model.RoleCoach),
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.announceRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.Announce{announce}, nil).Once()
			},
			wantLen: 1,
		},
		{
			name:   "athlete_signed",
			userID: athleteID,
			role:   string(model.RoleAthlete),
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(mem, nil).Once()
				m.announceRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.Announce{announce}, nil).Once()
				m.announceRepo.EXPECT().HasSignup(mock.Anything, annID, athleteID).Return(true, nil).Once()
			},
			wantLen:    1,
			wantSigned: true,
		},
		{
			name:   "athlete_not_member_empty",
			userID: athleteID,
			role:   string(model.RoleAthlete),
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(nil, model.ErrNotFound).Once()
			},
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := schedule_usecase.NewUseCase(m.announceRepo, m.clubRepo, m.membershipRepo)
			items, err := uc.List(context.Background(), tt.userID, tt.role)
			require.NoError(t, err)
			require.Len(t, items, tt.wantLen)
			if tt.wantLen > 0 {
				require.Equal(t, tt.wantSigned, items[0].SignedUp)
			}
		})
	}
}
