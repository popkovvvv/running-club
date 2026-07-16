//go:build unit

package schedule_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)

	tests := []struct {
		name         string
		req          dto.CreateAnnounceRequest
		before       func(m usecaseMocks, created **model.Announce)
		wantErr      error
		wantStartsOn string
	}{
		{
			name: "ok",
			req:  dto.CreateAnnounceRequest{Place: "ЛЭМЗ", Day: "Чт", Time: "19:50", Group: "Основная", Note: "test"},
			before: func(m usecaseMocks, created **model.Announce) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.announceRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Announce")).
					Run(func(_ context.Context, a *model.Announce) {
						*created = a
					}).
					Return(nil).Once()
				m.announceRepo.EXPECT().FindGoingAthletes(mock.Anything, mock.Anything).Return([]*model.User{}, nil).Once()
			},
		},
		{
			name: "ok_with_starts_on",
			req: dto.CreateAnnounceRequest{
				Place: "ЛЭМЗ", Day: "Чт", Time: "19:50", Group: "Основная", Note: "test", StartsOn: "2026-07-23",
			},
			before: func(m usecaseMocks, created **model.Announce) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.announceRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Announce")).
					Run(func(_ context.Context, a *model.Announce) {
						*created = a
					}).
					Return(nil).Once()
				m.announceRepo.EXPECT().FindGoingAthletes(mock.Anything, mock.Anything).Return([]*model.User{}, nil).Once()
			},
			wantStartsOn: "2026-07-23",
		},
		{
			name: "club_not_found",
			req:  dto.CreateAnnounceRequest{Place: "ЛЭМЗ", Day: "Чт", Time: "19:50", Group: "Основная", Note: "test"},
			before: func(m usecaseMocks, _ **model.Announce) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			var created *model.Announce
			if tt.before != nil {
				tt.before(m, &created)
			}
			uc := schedule_usecase.NewUseCase(m.announceRepo, m.clubRepo, m.membershipRepo, m.workoutRepo)
			view, err := uc.Publish(context.Background(), coachID, tt.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "ЛЭМЗ", view.Place)
			require.Equal(t, "Записаться", view.ScheduleCta)
			require.NotNil(t, created)
			if tt.wantStartsOn == "" {
				require.Nil(t, created.StartsOn)
				return
			}
			require.NotNil(t, created.StartsOn)
			require.Equal(t, tt.wantStartsOn, created.StartsOn.Format("2006-01-02"))
			require.Equal(t, time.UTC, created.StartsOn.Location())
		})
	}
}
