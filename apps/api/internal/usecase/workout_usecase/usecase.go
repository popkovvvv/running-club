package workout_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		workoutRepo workoutRepo
	}

	workoutRepo interface {
		Create(ctx context.Context, w *model.Workout) error
		FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error)
		FindOwnByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}
)

func NewUseCase(workoutRepo workoutRepo) *UseCase {
	return &UseCase{workoutRepo: workoutRepo}
}

var weekMeta = []struct{ Range, Plan string }{
	{"13.07 – 19.07", "25 км"},
	{"20.07 – 26.07", "27–28 км"},
	{"27.07 – 02.08", "30 км"},
	{"03.08 – 09.08", "33 км"},
}
