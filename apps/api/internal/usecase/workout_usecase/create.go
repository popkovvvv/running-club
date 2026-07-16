package workout_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (u *UseCase) Create(ctx context.Context, actorID uuid.UUID, role model.Role, req dto.CreateWorkoutRequest) (*dto.WorkoutView, error) {
	targetUserID := actorID
	var assignedBy *uuid.UUID
	var clubID *uuid.UUID

	if req.TargetUserID != nil && *req.TargetUserID != actorID {
		if role != model.RoleCoach {
			return nil, fmt.Errorf("forbidden assign: %w", model.ErrForbidden)
		}
		club, err := u.clubRepo.GetByCoachID(ctx, actorID)
		if err != nil {
			return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
		}
		membership, err := u.membershipRepo.GetByUserAndClub(ctx, *req.TargetUserID, club.ID)
		if err != nil {
			return nil, fmt.Errorf("membershipRepo.GetByUserAndClub: %w", err)
		}
		if membership.Status != model.MembershipActive {
			return nil, fmt.Errorf("student not active: %w", model.ErrForbidden)
		}
		targetUserID = *req.TargetUserID
		assignedBy = &actorID
		clubID = &club.ID
	} else if membership, err := u.membershipRepo.GetActiveByUser(ctx, actorID); err == nil {
		clubID = &membership.ClubID
	}

	w := buildWorkout(targetUserID, req, clubID, assignedBy, false)
	if w.Kind == model.WorkoutPlan {
		w.Status = model.WorkoutStatusPlanned
	}
	if err := u.workoutRepo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("workoutRepo.Create: %w", err)
	}
	v := mapWorkout(w)
	return &v, nil
}

func (u *UseCase) Get(ctx context.Context, actorID uuid.UUID, id uuid.UUID) (*dto.WorkoutView, error) {
	w, err := u.workoutRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.GetByID: %w", err)
	}
	if w.UserID != actorID {
		if _, err := u.clubRepo.GetByCoachID(ctx, actorID); err != nil {
			return nil, fmt.Errorf("access denied: %w", model.ErrForbidden)
		}
	}
	v := mapWorkout(w)
	return &v, nil
}

func (u *UseCase) Update(ctx context.Context, actorID uuid.UUID, id uuid.UUID, req dto.UpdateWorkoutRequest) (*dto.WorkoutView, error) {
	w, err := u.workoutRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.GetByID: %w", err)
	}
	if w.UserID != actorID {
		return nil, fmt.Errorf("access denied: %w", model.ErrForbidden)
	}
	if req.Status != nil {
		w.Status = model.WorkoutStatus(*req.Status)
	}
	if req.CompletedActivityID != nil {
		w.CompletedActivityID = req.CompletedActivityID
	}
	if err := u.workoutRepo.Update(ctx, w); err != nil {
		return nil, fmt.Errorf("workoutRepo.Update: %w", err)
	}
	v := mapWorkout(w)
	return &v, nil
}
