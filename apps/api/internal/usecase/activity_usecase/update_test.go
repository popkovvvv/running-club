//go:build unit

package activity_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateActivity(t *testing.T) {
	t.Parallel()
	ownerID := uuid.New()
	coachID := uuid.New()
	outsiderID := uuid.New()
	activityID := uuid.New()
	clubID := uuid.New()

	base := func() *model.Activity {
		a := model.NewActivity(ownerID, "Кросс", "Чт, 16.07.2026", 10, "", "", 0, 0, 0, "", 0, 0, 0, 0)
		a.ID = activityID
		a.Source = "manual"
		return a
	}

	title := "Индивидуальная"
	when := "2026-07-16"
	dist := 15.0
	duration := "1:05:00"
	pace := "4:20"
	hr := 148
	elev := 42.0

	tests := []struct {
		name    string
		actorID uuid.UUID
		role    model.Role
		req     dto.UpdateActivityRequest
		before  func(m usecaseMocks)
		wantErr error
		check   func(t *testing.T, view *dto.ActivityDetailView)
	}{
		{
			name:    "owner_ok",
			actorID: ownerID,
			role:    model.RoleAthlete,
			req: dto.UpdateActivityRequest{
				Title: &title, When: &when, DistKm: &dist,
				Duration: &duration, Pace: &pace, HR: &hr, ElevationGain: &elev,
			},
			before: func(m usecaseMocks) {
				a := base()
				m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(a, nil).Once()
				m.activityRepo.EXPECT().Update(mock.Anything, a).Return(nil).Once()
				m.workoutRepo.EXPECT().FindByCompletedActivity(mock.Anything, activityID).Return(nil, model.ErrNotFound).Once()
			},
			check: func(t *testing.T, view *dto.ActivityDetailView) {
				require.Equal(t, "Индивидуальная", view.Title)
				require.Equal(t, "15.0", view.Dist)
				require.Equal(t, "1:05:00", view.Time)
				require.Equal(t, "4:20", view.Pace)
				require.Equal(t, "148", view.HR)
				require.Equal(t, 42.0, view.Elevation)
				require.Equal(t, "Чт, 16.07.2026", view.When)
			},
		},
		{
			name:    "coach_of_member_ok",
			actorID: coachID,
			role:    model.RoleCoach,
			req:     dto.UpdateActivityRequest{Pace: &pace},
			before: func(m usecaseMocks) {
				a := base()
				m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(a, nil).Once()
				m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(&model.Club{ID: clubID}, nil).Once()
				m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, ownerID, clubID).Return(&model.Membership{
					UserID: ownerID, ClubID: clubID, Status: model.MembershipActive,
				}, nil).Once()
				m.activityRepo.EXPECT().Update(mock.Anything, a).Return(nil).Once()
				m.workoutRepo.EXPECT().FindByCompletedActivity(mock.Anything, activityID).Return(nil, model.ErrNotFound).Once()
			},
			check: func(t *testing.T, view *dto.ActivityDetailView) {
				require.Equal(t, "4:20", view.Pace)
			},
		},
		{
			name:    "outsider_forbidden",
			actorID: outsiderID,
			role:    model.RoleAthlete,
			req:     dto.UpdateActivityRequest{Pace: &pace},
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(base(), nil).Once()
			},
			wantErr: model.ErrForbidden,
		},
		{
			name:    "invalid_when",
			actorID: ownerID,
			role:    model.RoleAthlete,
			req:     dto.UpdateActivityRequest{When: ptr("16.07.2026")},
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(base(), nil).Once()
			},
			wantErr: model.ErrBadRequest,
		},
		{
			name:    "not_found",
			actorID: ownerID,
			role:    model.RoleAthlete,
			req:     dto.UpdateActivityRequest{Pace: &pace},
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(nil, model.ErrNotFound).Once()
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
			view, err := uc.Update(context.Background(), tt.actorID, tt.role, activityID, tt.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, view)
			if tt.check != nil {
				tt.check(t, view)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
