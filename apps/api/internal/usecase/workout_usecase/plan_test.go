//go:build unit

package workout_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlan(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	clubID := uuid.New()
	membership := model.MembershipFixture(uuid.New(), uid, clubID)
	planWorkout := model.NewWorkout(uid, model.WorkoutPlan, "Вт", "Интервалы", "Повторы", 8, "50 мин", "4:30", "160", 0)
	ownWorkout := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 6, "46 мин", "7:40", "140", 0)
	weeks := []*model.PlanWeek{
		model.NewPlanWeek(clubID, 0, "13.07 – 19.07", "25 км"),
		model.NewPlanWeek(clubID, 1, "20.07 – 26.07", "27–28 км"),
		model.NewPlanWeek(clubID, 2, "27.07 – 02.08", "30 км"),
		model.NewPlanWeek(clubID, 3, "03.08 – 09.08", "33 км"),
	}

	tests := []struct {
		name      string
		week      int
		before    func(m usecaseMocks)
		wantWeek  int
		wantDays  int
		wantMine  int
		wantRange string
		wantPlan  string
	}{
		{
			name: "ok_week0",
			week: 0,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{planWorkout}, nil).Once()
				m.workoutRepo.EXPECT().FindOwnByUser(mock.Anything, uid).
					Return([]*model.Workout{ownWorkout}, nil).Once()
			},
			wantWeek:  0,
			wantDays:  1,
			wantMine:  1,
			wantRange: "13.07 – 19.07",
			wantPlan:  "25 км",
		},
		{
			name: "clamp_negative",
			week: -1,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindOwnByUser(mock.Anything, uid).
					Return([]*model.Workout{}, nil).Once()
			},
			wantWeek:  0,
			wantDays:  0,
			wantMine:  0,
			wantRange: "13.07 – 19.07",
			wantPlan:  "25 км",
		},
		{
			name: "clamp_high",
			week: 99,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 3, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindOwnByUser(mock.Anything, uid).
					Return([]*model.Workout{}, nil).Once()
			},
			wantWeek:  3,
			wantDays:  0,
			wantMine:  0,
			wantRange: "03.08 – 09.08",
			wantPlan:  "33 км",
		},
		{
			name: "no_membership_empty_meta",
			week: 1,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 1, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindOwnByUser(mock.Anything, uid).
					Return([]*model.Workout{}, nil).Once()
			},
			wantWeek:  1,
			wantDays:  0,
			wantMine:  0,
			wantRange: "",
			wantPlan:  "",
		},
		{
			name: "no_plan_weeks_empty_meta",
			week: 0,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.PlanWeek{}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindOwnByUser(mock.Anything, uid).
					Return([]*model.Workout{}, nil).Once()
			},
			wantWeek:  0,
			wantDays:  0,
			wantMine:  0,
			wantRange: "",
			wantPlan:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo)
			res, err := uc.Plan(context.Background(), uid, tt.week)
			require.NoError(t, err)
			require.Equal(t, tt.wantWeek, res.WeekIndex)
			require.Equal(t, tt.wantRange, res.WeekRange)
			require.Equal(t, tt.wantPlan, res.WeekPlan)
			require.Len(t, res.Days, tt.wantDays)
			require.Len(t, res.Mine, tt.wantMine)
		})
	}
}
