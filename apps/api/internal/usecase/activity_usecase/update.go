package activity_usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Update(
	ctx context.Context,
	actorID uuid.UUID,
	role model.Role,
	activityID uuid.UUID,
	req dto.UpdateActivityRequest,
) (*dto.ActivityDetailView, error) {
	activity, err := u.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.GetByID: %w", err)
	}
	if err := u.ensureActivityAccess(ctx, actorID, role, activity.UserID); err != nil {
		return nil, err
	}

	if req.Title != nil {
		activity.Title = strings.TrimSpace(*req.Title)
	}
	if req.When != nil {
		d, err := time.Parse("2006-01-02", strings.TrimSpace(*req.When))
		if err != nil {
			return nil, fmt.Errorf("invalid when: %w", model.ErrBadRequest)
		}
		started := d.UTC()
		activity.StartedAt = &started
		activity.WhenLabel = formatWhenLabel(started)
	}
	if req.DistKm != nil {
		if *req.DistKm < 0 {
			return nil, fmt.Errorf("invalid distKm: %w", model.ErrBadRequest)
		}
		activity.DistKm = *req.DistKm
		activity.DistanceMeters = *req.DistKm * 1000
	}
	if req.Duration != nil {
		activity.Duration = strings.TrimSpace(*req.Duration)
		if sec, ok := parseDurationSeconds(activity.Duration); ok {
			activity.MovingSeconds = sec
			activity.ElapsedSeconds = sec
		}
	}
	if req.Pace != nil {
		activity.Pace = strings.TrimSpace(*req.Pace)
	}
	if req.HR != nil {
		if *req.HR < 0 {
			return nil, fmt.Errorf("invalid hr: %w", model.ErrBadRequest)
		}
		activity.HR = *req.HR
		activity.AverageHeartrate = *req.HR
	}
	if req.ElevationGain != nil {
		if *req.ElevationGain < 0 {
			return nil, fmt.Errorf("invalid elevationGain: %w", model.ErrBadRequest)
		}
		activity.ElevationGain = *req.ElevationGain
	}

	activity.UpdatedAt = time.Now().UTC()
	if err := u.activityRepo.Update(ctx, activity); err != nil {
		return nil, fmt.Errorf("activityRepo.Update: %w", err)
	}

	var linkedWorkoutID *uuid.UUID
	if w, err := u.workoutRepo.FindByCompletedActivity(ctx, activityID); err == nil {
		linkedWorkoutID = &w.ID
	}
	view := mapActivityDetail(activity, linkedWorkoutID)
	return &view, nil
}

func formatWhenLabel(t time.Time) string {
	days := []string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}
	return days[t.Weekday()] + ", " + t.Format("02.01.2006")
}

func parseDurationSeconds(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "—" {
		return 0, false
	}
	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}
	nums := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return 0, false
		}
		nums[i] = n
	}
	if len(nums) == 2 {
		return nums[0]*60 + nums[1], true
	}
	return nums[0]*3600 + nums[1]*60 + nums[2], true
}
