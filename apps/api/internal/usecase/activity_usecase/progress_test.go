//go:build unit

package activity_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProgress(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	now := time.Now().UTC()
	started := time.Date(now.Year(), time.June, 10, 12, 0, 0, 0, time.UTC)
	started2 := time.Date(now.Year(), time.July, 5, 12, 0, 0, 0, time.UTC)

	a1 := model.NewActivity(uid, "Run 1", "июнь", 100, "1:00", "5:30", 0, 0, 0, "", 0, 0, 0, 0)
	a1.StartedAt = &started
	a2 := model.NewActivity(uid, "Run 2", "июль", 50, "0:40", "5:45", 0, 0, 0, "", 0, 0, 0, 0)
	a2.StartedAt = &started2

	orphan := model.NewWorkout(uid, model.WorkoutOwn, "Пн", "Лёгкий", "Свой кросс", 6, "", 0)
	orphan.Status = model.WorkoutStatusCompleted
	orphanDate := time.Date(now.Year(), time.July, 12, 0, 0, 0, 0, time.UTC)
	orphan.ScheduledDate = &orphanDate

	races := []*model.Race{
		model.NewRace(uid, "Забег 1", "01.06", "10", "1:00", 0),
		model.NewRace(uid, "Забег 2", "01.07", "10", "1:00", 0),
	}

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
		check   func(t *testing.T, yearKm float64, yearTr, yearStarts int, monthsLen int)
	}{
		{
			name: "ok_from_activities_and_orphans",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindByUser(mock.Anything, uid).Return([]*model.Activity{a1, a2}, nil).Once()
				m.workoutRepo.EXPECT().FindCompletedWithoutActivity(mock.Anything, uid).Return([]*model.Workout{orphan}, nil).Once()
				m.activityRepo.EXPECT().FindRaces(mock.Anything, uid).Return(races, nil).Once()
				m.workoutRepo.EXPECT().FindByUser(mock.Anything, uid).Return([]*model.Workout{}, nil).Once()
			},
			check: func(t *testing.T, yearKm float64, yearTr, yearStarts int, monthsLen int) {
				require.Equal(t, 156.0, yearKm)
				require.Equal(t, 3, yearTr)
				require.Equal(t, 2, yearStarts)
				require.Equal(t, 2, monthsLen)
			},
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
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
			uc := newUC(m)
			res, err := uc.Progress(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.check(t, res.YearKm, res.YearTr, res.YearStarts, len(res.Months))
			require.Equal(t, "Июль", res.Months[0].M)
			require.Equal(t, "Июнь", res.Months[1].M)
		})
	}
}
