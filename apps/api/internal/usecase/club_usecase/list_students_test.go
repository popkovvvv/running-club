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

func TestListStudents(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", coachID)
	athlete := model.UserFixture(uuid.New(), "Иван Петров", "ivan@pulse.run", "hash", model.RoleAthlete)

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantLen int
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
				m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{athlete}, nil).Once()
			},
			wantLen: 1,
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
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo)
			students, err := uc.ListStudents(context.Background(), coachID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, students, tt.wantLen)
			require.Equal(t, "ИП", students[0].Init)
			require.Equal(t, "Иван Петров", students[0].Name)
		})
	}
}
