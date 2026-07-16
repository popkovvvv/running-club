//go:build unit

package workout_usecase_test

import (
	"context"
	"testing"
	"time"

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
	planWorkout := model.NewWorkout(uid, model.WorkoutPlan, "Вт", "Интервалы", "Повторы", 8, "160", 0)
	ownWorkout := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 6, "140", 0)
	weeks := []*model.PlanWeek{
		model.NewPlanWeek(clubID, 0, "13.07 – 19.07", "25 км"),
		model.NewPlanWeek(clubID, 1, "20.07 – 26.07", "27–28 км"),
		model.NewPlanWeek(clubID, 2, "27.07 – 02.08", "30 км"),
		model.NewPlanWeek(clubID, 3, "03.08 – 09.08", "33 км"),
	}
	expectWeekKm := func(m usecaseMocks) {
		m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, uid, mock.Anything).Return(12.5, nil).Once()
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
		wantKm    string
	}{
		{
			name: "ok_week0",
			week: 0,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{planWorkout}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutOwn).
					Return([]*model.Workout{ownWorkout}, nil).Once()
				expectWeekKm(m)
			},
			wantWeek:  0,
			wantDays:  1,
			wantMine:  1,
			wantRange: "13.07 – 19.07",
			wantPlan:  "25 км",
			wantKm:    "12.5",
		},
		{
			name: "clamp_negative",
			week: -1,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutOwn).
					Return([]*model.Workout{}, nil).Once()
				expectWeekKm(m)
			},
			wantWeek:  0,
			wantDays:  0,
			wantMine:  0,
			wantRange: "13.07 – 19.07",
			wantPlan:  "25 км",
			wantKm:    "12.5",
		},
		{
			name: "clamp_high",
			week: 99,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return(weeks, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 3, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 3, model.WorkoutOwn).
					Return([]*model.Workout{}, nil).Once()
				expectWeekKm(m)
			},
			wantWeek:  3,
			wantDays:  0,
			wantMine:  0,
			wantRange: "03.08 – 09.08",
			wantPlan:  "33 км",
			wantKm:    "12.5",
		},
		{
			name: "no_membership_default_range",
			week: 1,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 1, model.WorkoutPlan).
					Return([]*model.Workout{}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 1, model.WorkoutOwn).
					Return([]*model.Workout{}, nil).Once()
				expectWeekKm(m)
			},
			wantWeek:  1,
			wantDays:  0,
			wantMine:  0,
			wantRange: defaultRangeForWeek(1),
			wantPlan:  "",
			wantKm:    "12.5",
		},
		{
			name: "no_plan_weeks_default_range_and_volume",
			week: 0,
			before: func(m usecaseMocks) {
				m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(membership, nil).Once()
				m.planWeekRepo.EXPECT().FindByClub(mock.Anything, clubID).Return([]*model.PlanWeek{}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutPlan).
					Return([]*model.Workout{planWorkout}, nil).Once()
				m.workoutRepo.EXPECT().FindByUserWeek(mock.Anything, uid, 0, model.WorkoutOwn).
					Return([]*model.Workout{ownWorkout}, nil).Once()
				expectWeekKm(m)
			},
			wantWeek:  0,
			wantDays:  1,
			wantMine:  1,
			wantRange: defaultRangeForWeek(0),
			wantPlan:  "14 км",
			wantKm:    "12.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
			res, err := uc.Plan(context.Background(), uid, tt.week, 0, 0)
			require.NoError(t, err)
			require.Equal(t, tt.wantWeek, res.WeekIndex)
			require.Equal(t, tt.wantRange, res.WeekRange)
			require.Equal(t, tt.wantPlan, res.WeekPlan)
			require.Equal(t, tt.wantKm, res.WeekKm)
			require.Len(t, res.Days, tt.wantDays)
			require.Len(t, res.Mine, tt.wantMine)
		})
	}
}

func TestPlanMonth(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	inMonth := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	outMonth := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	planW := model.NewWorkout(uid, model.WorkoutPlan, "Чт", "Лёгкий", "План", 10, "", 0)
	planW.ScheduledDate = &inMonth
	ownW := model.NewWorkout(uid, model.WorkoutOwn, "Пт", "easy", "Свой", 7, "", 0)
	ownW.ScheduledDate = &inMonth
	other := model.NewWorkout(uid, model.WorkoutPlan, "Пн", "Лёгкий", "Июнь", 5, "", 0)
	other.ScheduledDate = &outMonth

	m := newMocks(t)
	m.workoutRepo.EXPECT().FindByUser(mock.Anything, uid).Return([]*model.Workout{planW, ownW, other}, nil).Once()
	m.membershipRepo.EXPECT().GetActiveByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
	m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, uid, mock.Anything).Return(3.0, nil).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	res, err := uc.Plan(context.Background(), uid, 0, 2026, 7)
	require.NoError(t, err)
	require.Len(t, res.Days, 1)
	require.Len(t, res.Mine, 1)
	require.Equal(t, "", res.WeekPlan)
	require.Equal(t, "3.0", res.WeekKm)
	require.NotEmpty(t, res.WeekRange)
}

func defaultRangeForWeek(weekIndex int) string {
	now := time.Now().UTC()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := now.AddDate(0, 0, -(weekday - 1))
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, weekIndex*7)
	end := start.AddDate(0, 0, 6)
	return start.Format("02.01") + " – " + end.Format("02.01")
}
