package workout_usecase

import (
	"github.com/google/uuid"
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
	tag := w.Tag
	if tag == "" {
		tag = workoutTypeLabel(w.WorkoutType)
	}
	return dto.NewWorkoutView(
		w.ID, string(w.Kind), string(w.WorkoutType), w.DayLabel, tag, w.Title, w.Description,
		w.DistKm, w.Duration, w.Pace, w.HR, w.WeekIndex,
		dto.FormatDate(w.ScheduledDate), string(w.Status),
		w.CompletedActivityID, w.AssignedBy, w.IsClubTemplate, segments,
	)
}

func workoutTypeLabel(t model.WorkoutType) string {
	labels := map[model.WorkoutType]string{
		model.WorkoutTypeEasy:     "Лёгкий",
		model.WorkoutTypeLong:     "Длинный",
		model.WorkoutTypeTempo:    "Темповый",
		model.WorkoutTypeInterval: "Интервалы",
		model.WorkoutTypeFartlek:  "Фартlek",
		model.WorkoutTypeRecovery: "Восстановление",
		model.WorkoutTypeHills:    "Горки",
		model.WorkoutTypeRace:     "Старт",
		model.WorkoutTypeCross:    "Кросс",
		model.WorkoutTypeRest:     "Отдых",
	}
	if l, ok := labels[t]; ok {
		return l
	}
	return "Тренировка"
}

func buildWorkout(targetUserID uuid.UUID, req dto.CreateWorkoutRequest, clubID *uuid.UUID, assignedBy *uuid.UUID, isTemplate bool) *model.Workout {
	kind := model.WorkoutKind(req.Kind)
	if kind == "" {
		kind = model.WorkoutOwn
	}
	workoutType := model.WorkoutTypeEasy
	if req.WorkoutType != "" && model.ValidWorkoutType(req.WorkoutType) {
		workoutType = model.WorkoutType(req.WorkoutType)
	}
	distKm := req.DistKm
	if len(req.Segments) > 0 {
		distKm = 0
	}
	w := model.NewWorkout(targetUserID, kind, req.DayLabel, req.Tag, req.Title, distKm, req.Duration, req.Pace, req.HR, req.WeekIndex)
	w.ClubID = clubID
	w.WorkoutType = workoutType
	w.Description = req.Description
	w.AssignedBy = assignedBy
	w.IsClubTemplate = isTemplate
	if req.ScheduledDate != nil {
		if d, err := parseDate(*req.ScheduledDate); err == nil {
			w.ScheduledDate = &d
		}
	}
	for i, s := range req.Segments {
		w.AddSegment(s.Kind, s.Title, s.DistKm, s.Pace, i)
	}
	return w
}
