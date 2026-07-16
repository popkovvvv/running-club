package activity_usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) GetByID(ctx context.Context, actorID uuid.UUID, role model.Role, activityID uuid.UUID) (*dto.ActivityDetailView, error) {
	activity, err := u.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.GetByID: %w", err)
	}
	if err := u.ensureActivityAccess(ctx, actorID, role, activity.UserID); err != nil {
		return nil, err
	}
	var linkedWorkoutID *uuid.UUID
	if w, err := u.workoutRepo.FindByCompletedActivity(ctx, activityID); err == nil {
		linkedWorkoutID = &w.ID
	}
	view := mapActivityDetail(activity, linkedWorkoutID)
	return &view, nil
}

func (u *UseCase) GetStreams(ctx context.Context, actorID uuid.UUID, role model.Role, activityID uuid.UUID) ([]dto.ActivityStreamView, error) {
	activity, err := u.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.GetByID: %w", err)
	}
	if err := u.ensureActivityAccess(ctx, actorID, role, activity.UserID); err != nil {
		return nil, err
	}
	streams, err := u.activityStreamRepo.FindByActivityID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("activityStreamRepo.FindByActivityID: %w", err)
	}
	out := make([]dto.ActivityStreamView, 0, len(streams))
	for _, s := range streams {
		out = append(out, dto.NewActivityStreamView(s.Type, s.DataJSON))
	}
	return out, nil
}

func (u *UseCase) ensureActivityAccess(ctx context.Context, actorID uuid.UUID, role model.Role, ownerID uuid.UUID) error {
	if actorID == ownerID {
		return nil
	}
	if role != model.RoleCoach {
		return fmt.Errorf("access denied: %w", model.ErrForbidden)
	}
	club, err := u.clubRepo.GetByCoachID(ctx, actorID)
	if err != nil {
		return fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	membership, err := u.membershipRepo.GetByUserAndClub(ctx, ownerID, club.ID)
	if err != nil {
		return fmt.Errorf("membershipRepo.GetByUserAndClub: %w", err)
	}
	if membership.Status != model.MembershipActive {
		return fmt.Errorf("not club member: %w", model.ErrForbidden)
	}
	return nil
}

func mapActivityDetail(a *model.Activity, linkedWorkoutID *uuid.UUID) dto.ActivityDetailView {
	return dto.NewActivityDetailView(
		a.ID, a.Title, a.WhenLabel, a.StartedAt,
		formatKm(a.DistKm), a.Duration, a.Pace, strconv.Itoa(a.HR),
		a.MaxHeartrate, a.MovingSeconds, a.ElapsedSeconds,
		a.Kudos, a.Comments, a.RouteSVG, a.Polyline,
		a.StartX, a.StartY, a.EndX, a.EndY,
		a.Source, a.SportType, a.ElevationGain, a.Visibility, a.ExternalID,
		linkedWorkoutID,
	)
}
