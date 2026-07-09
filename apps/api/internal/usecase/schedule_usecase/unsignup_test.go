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

func TestUnsignup(t *testing.T) {
	t.Parallel()
	clubID := uuid.New()
	annID := uuid.New()
	athleteID := uuid.New()
	announce := model.AnnounceFixture(annID, clubID, "Зина", "Вт", "19:50", "Основная")
	announce.GoingCount = 1
	announceAfter := model.AnnounceFixture(annID, clubID, "Зина", "Вт", "19:50", "Основная")
	announceAfter.GoingCount = 0

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announce, nil).Once()
				m.announceRepo.EXPECT().DeleteSignup(mock.Anything, annID, athleteID).Return(nil).Once()
				m.announceRepo.EXPECT().IncGoing(mock.Anything, annID, -1).Return(nil).Once()
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announceAfter, nil).Once()
			},
		},
		{
			name: "not_signed_up",
			before: func(m usecaseMocks) {
				m.announceRepo.EXPECT().GetByID(mock.Anything, annID).Return(announce, nil).Once()
				m.announceRepo.EXPECT().DeleteSignup(mock.Anything, annID, athleteID).Return(model.ErrNotFound).Once()
			},
			wantErr: model.ErrNotSignedUp,
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
			view, err := uc.Unsignup(context.Background(), athleteID, annID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "Записаться", view.ScheduleCta)
			require.False(t, view.SignedUp)
		})
	}
}
