package workout_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Create(ctx context.Context, userID uuid.UUID, req dto.CreateWorkoutRequest) (*dto.WorkoutView, error) {
	kind := model.WorkoutKind(req.Kind)
	if kind == "" {
		kind = model.WorkoutOwn
	}
	w := model.NewWorkout(userID, kind, req.DayLabel, req.Tag, req.Title, req.DistKm, req.Duration, req.Pace, req.HR, req.WeekIndex)
	for i, s := range req.Segments {
		w.AddSegment(s.Kind, s.Title, s.DistKm, s.Pace, i)
	}
	if err := u.workoutRepo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("workoutRepo.Create: %w", err)
	}
	v := mapWorkout(w)
	return &v, nil
}
