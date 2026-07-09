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

func TestListStudents(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	athleteID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	athlete := model.UserFixture(athleteID, "Иван Петров", "ivan@pulse.run", "hash", model.RoleAthlete)
	planWeek := model.NewPlanWeek(clubID, 0, "13.07 – 19.07", "25 км")

	tests := []struct {
		name     string
		before   func(m usecaseMocks)
		wantLen  int
		wantKm   string
		wantSub  string
		wantComp int
		wantErr  error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.planWeekRepo.EXPECT().GetByClubAndIndex(mock.Anything, clubID, 0).Return(planWeek, nil).Once()
				m.activityRepo.EXPECT().SumDistByUser(mock.Anything, athleteID).Return(24.6, nil).Once()
				m.announceRepo.EXPECT().NextLabelForAthlete(mock.Anything, clubID, athleteID).Return("Вт, 21 июля · Стадион «Зина»", nil).Once()
			},
			wantLen:  1,
			wantKm:   "24.6",
			wantSub:  "Вт, 21 июля · Стадион «Зина»",
			wantComp: 98,
		},
		{
			name: "no_upcoming_announce",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.planWeekRepo.EXPECT().GetByClubAndIndex(mock.Anything, clubID, 0).Return(nil, model.ErrNotFound).Once()
				m.activityRepo.EXPECT().SumDistByUser(mock.Anything, athleteID).Return(0, nil).Once()
				m.announceRepo.EXPECT().NextLabelForAthlete(mock.Anything, clubID, athleteID).Return("", model.ErrNotFound).Once()
			},
			wantLen:  1,
			wantKm:   "0.0",
			wantSub:  "Нет записи",
			wantComp: 0,
		},
		{
			name: "club_not_found",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(nil, model.ErrNotFound).Once()
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
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo, m.activityRepo, m.announceRepo, m.planWeekRepo)
			students, err := uc.ListStudents(context.Background(), coachID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, students, tt.wantLen)
			require.Equal(t, "ИП", students[0].Init)
			require.Equal(t, "Иван Петров", students[0].Name)
			require.Equal(t, tt.wantKm, students[0].Km)
			require.Equal(t, tt.wantSub, students[0].Sub)
			require.Equal(t, tt.wantComp, students[0].Comp)
		})
	}
}
