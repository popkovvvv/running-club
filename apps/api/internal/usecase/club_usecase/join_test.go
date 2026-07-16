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

func TestJoin(t *testing.T) {
	t.Parallel()
	clubID := uuid.New()
	club := model.ClubFixture(clubID, "PULSE", "PULSE-7K42", "#ff5c22", uuid.Nil)

	tests := []struct {
		name    string
		code    string
		before  func(m usecaseMocks, uid uuid.UUID)
		wantErr error
	}{
		{
			name: "ok",
			code: "PULSE-7K42",
			before: func(m usecaseMocks, uid uuid.UUID) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
				m.clubRepo.EXPECT().GetByInviteCode(mock.Anything, "PULSE-7K42").Return(club, nil).Once()
				m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, uid, clubID).Return(nil, model.ErrNotFound).Once()
				m.membershipRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Membership")).Return(nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, clubID).Return(1, nil).Once()
			},
		},
		{
			name: "ok_reactivate",
			code: "PULSE-7K42",
			before: func(m usecaseMocks, uid uuid.UUID) {
				existing := model.MembershipFixture(uuid.New(), uid, clubID)
				existing.Status = model.MembershipLeft
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
				m.clubRepo.EXPECT().GetByInviteCode(mock.Anything, "PULSE-7K42").Return(club, nil).Once()
				m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, uid, clubID).Return(existing, nil).Once()
				m.membershipRepo.EXPECT().UpdateStatus(mock.Anything, existing.ID, model.MembershipActive).Return(nil).Once()
				m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, clubID).Return(1, nil).Once()
			},
		},
		{
			name: "bad_code",
			code: "NOPE",
			before: func(m usecaseMocks, uid uuid.UUID) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
				m.clubRepo.EXPECT().GetByInviteCode(mock.Anything, "NOPE").Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrInvalidInviteCode,
		},
		{
			name:    "empty_code",
			code:    "",
			wantErr: model.ErrInvalidInviteCode,
		},
		{
			name: "already_member",
			code: "PULSE-7K42",
			before: func(m usecaseMocks, uid uuid.UUID) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).
					Return(model.NewMembership(uid, clubID), nil).Once()
			},
			wantErr: model.ErrAlreadyMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uid := uuid.New()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m, uid)
			}
			uc := club_usecase.NewUseCase(m.clubRepo, m.membershipRepo, m.userRepo, m.activityRepo, m.announceRepo, m.planWeekRepo, m.workoutRepo)
			view, err := uc.Join(context.Background(), uid, tt.code)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, clubID, view.ID)
			require.Equal(t, "PULSE", view.Name)
		})
	}
}
