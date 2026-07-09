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

func TestSignup(t *testing.T) {
	t.Parallel()
	clubID := uuid.New()
	annID := uuid.New()
	athleteID := uuid.New()
	announce := model.AnnounceFixture(annID, clubID, "Зина", "Вт", "19:50", "Основная")
	announceAfter := model.AnnounceFixture(annID, clubID, "Зина", "Вт", "19:50", "Основная")
	announceAfter.GoingCount = 1

	tests := []struct {
		name      string
		before    func(m usecaseMocks)
		wantErr   error
		wantCTA   string
		wantGoing int
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announce, nil).Once()
				m.announceRepo.EXPECT().HasSignup(mock.Anything, annID, athleteID).Return(false, nil).Once()
				m.announceRepo.EXPECT().CreateSignup(mock.Anything, mock.AnythingOfType("*model.AnnounceSignup")).Return(nil).Once()
				m.announceRepo.EXPECT().IncGoing(mock.Anything, annID, 1).Return(nil).Once()
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announceAfter, nil).Once()
			},
			wantCTA:   "Вы записаны",
			wantGoing: 1,
		},
		{
			name: "already_signed_up",
			before: func(m usecaseMocks) {
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announce, nil).Once()
				m.announceRepo.EXPECT().HasSignup(mock.Anything, annID, athleteID).Return(true, nil).Once()
			},
			wantErr: model.ErrAlreadySignedUp,
		},
		{
			name: "announce_not_found",
			before: func(m usecaseMocks) {
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(nil, model.ErrNotFound).Once()
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
			view, err := uc.Signup(context.Background(), athleteID, annID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantCTA, view.ScheduleCta)
			require.True(t, view.SignedUp)
			require.Equal(t, tt.wantGoing, view.Going)
		})
	}
}
