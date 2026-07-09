//go:build unit

package schedule_usecase_test

import (
	"context"
	"testing"

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
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.announceRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Announce")).Return(nil).Once()
			},
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
			uc := schedule_usecase.NewUseCase(m.announceRepo, m.clubRepo, m.membershipRepo)
			view, err := uc.Publish(context.Background(), coachID, dto.CreateAnnounceRequest{
				Place: "ЛЭМЗ", Day: "Чт", Time: "19:50", Group: "Основная", Note: "test",
			})
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "ЛЭМЗ", view.Place)
			require.Equal(t, "Записаться", view.ScheduleCta)
		})
	}
}
