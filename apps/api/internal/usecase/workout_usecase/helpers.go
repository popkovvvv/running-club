package workout_usecase

import (
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) SegmentTotal(segments []dto.SegmentInput) float64 {
	var sum float64
	for _, s := range segments {
		sum += s.DistKm
	}
	return sum
}

func mapWorkouts(items []*model.Workout) []dto.WorkoutView {
	out := make([]dto.WorkoutView, 0, len(items))
	for _, w := range items {
		out = append(out, mapWorkout(w))
	}
	return out
}

func mapWorkout(w *model.Workout) dto.WorkoutView {
	segments := make([]dto.SegmentView, 0, len(w.Segments))
	for _, s := range w.Segments {
		segments = append(segments, dto.NewSegmentView(s.ID, s.Kind, s.Title, s.DistKm, s.Pace))
	}
	return dto.NewWorkoutView(
		w.ID, string(w.Kind), w.DayLabel, w.Tag, w.Title,
		w.DistKm, w.Duration, w.Pace, w.HR, segments,
	)
}
