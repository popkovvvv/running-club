//go:build unit

package activity_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProgress(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	months := []model.MonthStat{
		model.NewMonthStat("Июн", 100, 10, "5:30", "+10"),
		model.NewMonthStat("Июл", 50, 5, "5:45", "+5"),
	}
	races := []*model.Race{
		model.NewRace(uid, "Забег 1", "01.06", "10", "1:00", 0),
		model.NewRace(uid, "Забег 2", "01.07", "10", "1:00", 0),
	}

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindMonthStats(mock.Anything, uid).Return(months, nil).Once()
				m.activityRepo.EXPECT().FindRaces(mock.Anything, uid).Return(races, nil).Once()
			},
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindMonthStats(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
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
			uc := activity_usecase.NewUseCase(m.activityRepo)
			res, err := uc.Progress(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 150.0, res.YearKm)
			require.Equal(t, 15, res.YearTr)
			require.Equal(t, 2, res.YearStarts)
			require.Len(t, res.Months, 2)
			require.Equal(t, "Июн", res.Months[0].M)
			require.Equal(t, "Июл", res.Months[1].M)
		})
	}
}
