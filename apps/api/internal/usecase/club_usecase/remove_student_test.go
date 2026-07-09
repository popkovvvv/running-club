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

func TestRemoveStudent(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	studentID := uuid.New()
	clubID := uuid.New()
	memID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	mem := model.MembershipFixture(memID, studentID, clubID)

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, studentID, clubID).Return(mem, nil).Once()
				m.membershipRepo.EXPECT().UpdateStatus(mock.Anything, memID, model.MembershipRemoved).Return(nil).Once()
			},
		},
		{
			name: "membership_not_found",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, studentID, clubID).Return(nil, model.ErrNotFound).Once()
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
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo)
			err := uc.RemoveStudent(context.Background(), coachID, studentID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
