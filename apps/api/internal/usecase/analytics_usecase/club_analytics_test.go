//go:build unit

package analytics_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/analytics_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClubAnalytics(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	athleteID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	athlete := model.UserFixture(athleteID, "Иван Петров", "ivan@pulse.run", "hash", model.RoleAthlete)
	planWeek := model.NewPlanWeek(clubID, 0, "13.07 – 19.07", "25 км")

	tests := []struct {
		name          string
		before        func(m usecaseMocks)
		wantLen       int
		wantClubKm    float64
		wantStudentKm string
		wantComp      int
		wantErr       error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.activityRepo.EXPECT().SumDistByClubAthletes(mock.Anything, clubID).Return(48.2, nil).Once()
				m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, athleteID, mock.Anything).Return(24.6, nil).Once()
				m.workoutRepo.EXPECT().SumPlanDistByUserWeek(mock.Anything, athleteID, 0).Return(25.0, nil).Once()
			},
			wantLen:       1,
			wantClubKm:    48.2,
			wantStudentKm: "24.6",
			wantComp:      98,
		},
		{
			name: "no_plan",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.activityRepo.EXPECT().SumDistByClubAthletes(mock.Anything, clubID).Return(10.0, nil).Once()
				m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, athleteID, mock.Anything).Return(10.0, nil).Once()
				m.workoutRepo.EXPECT().SumPlanDistByUserWeek(mock.Anything, athleteID, 0).Return(0, nil).Once()
				m.planWeekRepo.EXPECT().GetByClubAndIndex(mock.Anything, clubID, 0).Return(nil, model.ErrNotFound).Once()
			},
			wantLen:       1,
			wantClubKm:    10.0,
			wantStudentKm: "10.0",
			wantComp:      0,
		},
		{
			name: "fallback_plan_label",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.activityRepo.EXPECT().SumDistByClubAthletes(mock.Anything, clubID).Return(12.5, nil).Once()
				m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, athleteID, mock.Anything).Return(12.5, nil).Once()
				m.workoutRepo.EXPECT().SumPlanDistByUserWeek(mock.Anything, athleteID, 0).Return(0, nil).Once()
				m.planWeekRepo.EXPECT().GetByClubAndIndex(mock.Anything, clubID, 0).Return(planWeek, nil).Once()
			},
			wantLen:       1,
			wantClubKm:    12.5,
			wantStudentKm: "12.5",
			wantComp:      50,
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
			uc := analytics_usecase.NewUseCase(m.clubRepo, m.userRepo, m.activityRepo, m.planWeekRepo, m.workoutRepo)
			res, err := uc.ClubAnalytics(context.Background(), coachID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantClubKm, res.ClubKm)
			require.Len(t, res.Students, tt.wantLen)
			require.Equal(t, "ИП", res.Students[0].Init)
			require.Equal(t, tt.wantStudentKm, res.Students[0].Km)
			require.Equal(t, tt.wantComp, res.Students[0].Comp)
		})
	}
}
