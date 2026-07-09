//go:build unit

package workout_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	uid := uuid.New()

	tests := []struct {
		name   string
		req    dto.CreateWorkoutRequest
		want   float64
		before func(m usecaseMocks)
	}{
		{
			name: "own_easy",
			req:  dto.CreateWorkoutRequest{Kind: "own", DayLabel: "Ср", Tag: "Лёгкий", Title: "Кросс", DistKm: 6, Duration: "~46 мин", Pace: "7:40"},
			want: 6,
			before: func(m usecaseMocks) {
				m.workoutRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Workout")).Return(nil).Once()
			},
		},
		{
			name: "builder_segments",
			req: dto.CreateWorkoutRequest{Kind: "builder", Title: "ОФП", Segments: []dto.SegmentInput{
				{Kind: "Разминка", Title: "Лёгкий", DistKm: 2, Pace: "7:40"},
				{Kind: "Основная", Title: "ОФП", DistKm: 5, Pace: "смеш."},
				{Kind: "Заминка", Title: "Легко", DistKm: 1, Pace: "8:00"},
			}},
			want: 8,
			before: func(m usecaseMocks) {
				m.workoutRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Workout")).Return(nil).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := workout_usecase.NewUseCase(m.workoutRepo)
			view, err := uc.Create(context.Background(), uid, tt.req)
			require.NoError(t, err)
			require.InDelta(t, tt.want, view.DistKm, 0.01)
		})
	}
}

func TestSegmentTotal(t *testing.T) {
	t.Parallel()
	m := newMocks(t)
	uc := workout_usecase.NewUseCase(m.workoutRepo)
	require.Equal(t, 8.0, uc.SegmentTotal([]dto.SegmentInput{{DistKm: 2}, {DistKm: 5}, {DistKm: 1}}))
}
