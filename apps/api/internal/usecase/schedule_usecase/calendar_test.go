//go:build unit

package schedule_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCalendar(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	athleteID := uuid.New()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	mem := model.MembershipFixture(uuid.New(), athleteID, clubID)

	now := time.Now()
	sessionDay := time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, time.UTC)
	otherMonth := time.Date(now.Year(), now.Month()+1, 3, 0, 0, 0, 0, time.UTC)
	annWithSession := model.AnnounceFixture(uuid.New(), clubID, "Зина", "Вт", "19:50", "Основная")
	annWithSession.StartsOn = &sessionDay
	annOtherMonth := model.AnnounceFixture(uuid.New(), clubID, "ЛЭМЗ", "Чт", "19:50", "Основная")
	annOtherMonth.StartsOn = &otherMonth
	annNoDate := model.AnnounceFixture(uuid.New(), clubID, "Парк", "Сб", "10:00", "Основная")

	tests := []struct {
		name    string
		userID  uuid.UUID
		role    string
		before  func(m usecaseMocks)
		wantHas int
		empty   bool
	}{
		{
			name:   "coach_marks_current_month_days",
			userID: coachID,
			role:   string(model.RoleCoach),
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.announceRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.Announce{
					annWithSession, annOtherMonth, annNoDate,
				}, nil).Once()
			},
			wantHas: 1,
		},
		{
			name:   "athlete_loads_club_accent",
			userID: athleteID,
			role:   string(model.RoleAthlete),
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(mem, nil).Once()
				m.clubRepo.EXPECT().GetByID(mock.Anything, clubID).Return(club, nil).Once()
				m.announceRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.Announce{annWithSession}, nil).Once()
			},
			wantHas: 1,
		},
		{
			name:   "athlete_not_member_empty",
			userID: athleteID,
			role:   string(model.RoleAthlete),
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, athleteID).Return(nil, model.ErrNotFound).Once()
			},
			empty: true,
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
			res, err := uc.Calendar(context.Background(), tt.userID, tt.role)
			require.NoError(t, err)
			if tt.empty {
				require.Empty(t, res.Cells)
				return
			}
			require.NotEmpty(t, res.Cells)
			var blanks, hasCount, todayCount int
			for _, c := range res.Cells {
				if c.Blank {
					blanks++
					require.Equal(t, 0, c.N)
					require.Equal(t, "transparent", c.Dot)
					require.Equal(t, "transparent", c.Bg)
					continue
				}
				if c.Has {
					hasCount++
				}
				if c.IsToday {
					todayCount++
					require.Equal(t, "#ff5c22", c.Bg)
					require.Equal(t, "#ffffff", c.Fg)
				} else if c.Has {
					require.Equal(t, "rgba(255,92,34,.15)", c.Bg)
					require.Equal(t, "#ff5c22", c.Dot)
				}
			}
			require.Greater(t, blanks, 0)
			require.Equal(t, tt.wantHas, hasCount)
			require.Equal(t, 1, todayCount)
		})
	}
}
