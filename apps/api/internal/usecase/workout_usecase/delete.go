package workout_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.workoutRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("workoutRepo.Delete: %w", err)
	}
	return nil
}
