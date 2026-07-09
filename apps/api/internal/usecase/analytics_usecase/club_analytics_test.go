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

	tests := []struct {
		name           string
		before         func(m usecaseMocks)
		wantLen        int
		wantClubKm     float64
		wantAttendance int
		wantStudentKm  string
		wantErr        error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.activityRepo.EXPECT().SumDistByClubAthletes(mock.Anything, clubID).Return(48.2, nil).Once()
				m.announceRepo.EXPECT().AttendanceStats(mock.Anything, clubID).Return(12, 14, nil).Once()
				m.activityRepo.EXPECT().SumDistByUser(mock.Anything, athleteID).Return(24.6, nil).Once()
			},
			wantLen:        1,
			wantClubKm:     48.2,
			wantAttendance: 86,
			wantStudentKm:  "24.6",
		},
		{
			name: "no_announces",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
				m.activityRepo.EXPECT().SumDistByClubAthletes(mock.Anything, clubID).Return(10.0, nil).Once()
				m.announceRepo.EXPECT().AttendanceStats(mock.Anything, clubID).Return(0, 0, nil).Once()
				m.activityRepo.EXPECT().SumDistByUser(mock.Anything, athleteID).Return(10.0, nil).Once()
			},
			wantLen:        1,
			wantClubKm:     10.0,
			wantAttendance: 0,
			wantStudentKm:  "10.0",
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
			uc := analytics_usecase.NewUseCase(m.clubRepo, m.userRepo, m.activityRepo, m.announceRepo)
			res, err := uc.ClubAnalytics(context.Background(), coachID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantClubKm, res.ClubKm)
			require.Equal(t, tt.wantAttendance, res.Attendance)
			require.Len(t, res.Students, tt.wantLen)
			require.Equal(t, "ИП", res.Students[0].Init)
			require.Equal(t, tt.wantStudentKm, res.Students[0].Km)
		})
	}
}
