package plan_usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func buildTemplateWorkout(coachID, clubID uuid.UUID, weekIndex int, req dto.CreateWorkoutRequest) *model.Workout {
	workoutType := model.WorkoutTypeEasy
	if req.WorkoutType != "" && model.ValidWorkoutType(req.WorkoutType) {
		workoutType = model.WorkoutType(req.WorkoutType)
	}
	distKm := req.DistKm
	if len(req.Segments) > 0 {
		distKm = 0
	}
	w := &model.Workout{
		ID:             uuid.New(),
		ClubID:         &clubID,
		UserID:         coachID,
		Kind:           model.WorkoutPlan,
		WorkoutType:    workoutType,
		DayLabel:       req.DayLabel,
		Tag:            req.Tag,
		Title:          req.Title,
		Description:    req.Description,
		DistKm:         distKm,
		Duration:       req.Duration,
		Pace:           req.Pace,
		HR:             req.HR,
		WeekIndex:      weekIndex,
		Status:         model.WorkoutStatusPlanned,
		IsClubTemplate: true,
		CreatedAt:      time.Now().UTC(),
	}
	if req.ScheduledDate != nil {
		if d, err := time.Parse("2006-01-02", *req.ScheduledDate); err == nil {
			w.ScheduledDate = &d
		}
	}
	for i, s := range req.Segments {
		w.Segments = append(w.Segments, model.NewSegment(s.Kind, s.Title, s.DistKm, s.Pace, i))
	}
	return w
}

func mapWorkouts(items []*model.Workout) []dto.WorkoutView {
	out := make([]dto.WorkoutView, 0, len(items))
	for _, w := range items {
		segments := make([]dto.SegmentView, 0, len(w.Segments))
		for _, s := range w.Segments {
			segments = append(segments, dto.NewSegmentView(s.ID, s.Kind, s.Title, s.DistKm, s.Pace))
		}
		tag := w.Tag
		if tag == "" {
			tag = string(w.WorkoutType)
		}
		out = append(out, dto.NewWorkoutView(
			w.ID, string(w.Kind), string(w.WorkoutType), w.DayLabel, tag, w.Title, w.Description,
			w.DistKm, w.Duration, w.Pace, w.HR, w.WeekIndex,
			dto.FormatDate(w.ScheduledDate), string(w.Status),
			w.CompletedActivityID, w.AssignedBy, w.IsClubTemplate, segments,
		))
	}
	return out
}
