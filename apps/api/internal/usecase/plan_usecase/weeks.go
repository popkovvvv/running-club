package plan_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) ListWeeks(ctx context.Context, coachID uuid.UUID) ([]dto.PlanWeekView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	weeks, err := u.planWeekRepo.FindByClub(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("planWeekRepo.FindByClub: %w", err)
	}
	out := make([]dto.PlanWeekView, 0, len(weeks))
	for _, w := range weeks {
		out = append(out, dto.PlanWeekView{
			WeekIndex:  w.WeekIndex,
			RangeLabel: w.RangeLabel,
			PlanLabel:  w.PlanLabel,
		})
	}
	return out, nil
}

func (u *UseCase) UpsertWeek(ctx context.Context, coachID uuid.UUID, weekIndex int, req dto.UpsertPlanWeekRequest) (*dto.PlanWeekView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	pw := model.NewPlanWeek(club.ID, weekIndex, req.RangeLabel, req.PlanLabel)
	if err := u.planWeekRepo.Upsert(ctx, pw); err != nil {
		return nil, fmt.Errorf("planWeekRepo.Upsert: %w", err)
	}
	return &dto.PlanWeekView{
		WeekIndex:  pw.WeekIndex,
		RangeLabel: pw.RangeLabel,
		PlanLabel:  pw.PlanLabel,
	}, nil
}

func (u *UseCase) GetTemplate(ctx context.Context, coachID uuid.UUID, weekIndex int) (*dto.TemplateResponse, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	workouts, err := u.workoutRepo.FindClubTemplates(ctx, club.ID, weekIndex)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindClubTemplates: %w", err)
	}
	return &dto.TemplateResponse{
		WeekIndex: weekIndex,
		Workouts:  mapWorkouts(workouts),
	}, nil
}

func (u *UseCase) SaveTemplate(ctx context.Context, coachID uuid.UUID, weekIndex int, req dto.SaveTemplateRequest) (*dto.TemplateResponse, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	workouts := make([]*model.Workout, 0, len(req.Workouts))
	for _, item := range req.Workouts {
		w := buildTemplateWorkout(coachID, club.ID, weekIndex, item)
		workouts = append(workouts, w)
	}
	if err := u.workoutRepo.ReplaceClubTemplates(ctx, club.ID, weekIndex, workouts); err != nil {
		return nil, fmt.Errorf("workoutRepo.ReplaceClubTemplates: %w", err)
	}
	saved, err := u.workoutRepo.FindClubTemplates(ctx, club.ID, weekIndex)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindClubTemplates: %w", err)
	}
	return &dto.TemplateResponse{
		WeekIndex: weekIndex,
		Workouts:  mapWorkouts(saved),
	}, nil
}

func (u *UseCase) Publish(ctx context.Context, coachID uuid.UUID, weekIndex int) error {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	templates, err := u.workoutRepo.FindClubTemplates(ctx, club.ID, weekIndex)
	if err != nil {
		return fmt.Errorf("workoutRepo.FindClubTemplates: %w", err)
	}
	athletes, err := u.userRepo.FindAthletesByClub(ctx, club.ID)
	if err != nil {
		return fmt.Errorf("userRepo.FindAthletesByClub: %w", err)
	}
	for _, athlete := range athletes {
		if err := u.workoutRepo.DeleteClubAssignedPlans(ctx, athlete.ID, weekIndex); err != nil {
			return fmt.Errorf("workoutRepo.DeleteClubAssignedPlans: %w", err)
		}
		for _, tmpl := range templates {
			copy := cloneForAthlete(tmpl, athlete.ID, club.ID)
			if err := u.workoutRepo.Create(ctx, copy); err != nil {
				return fmt.Errorf("workoutRepo.Create: %w", err)
			}
		}
	}
	return nil
}

func cloneForAthlete(src *model.Workout, athleteID, clubID uuid.UUID) *model.Workout {
	w := &model.Workout{
		ID:            uuid.New(),
		ClubID:        &clubID,
		UserID:        athleteID,
		Kind:          model.WorkoutPlan,
		WorkoutType:   src.WorkoutType,
		DayLabel:      src.DayLabel,
		Tag:           src.Tag,
		Title:         src.Title,
		Description:   src.Description,
		DistKm:        src.DistKm,
		Duration:      src.Duration,
		Pace:          src.Pace,
		HR:            src.HR,
		WeekIndex:     src.WeekIndex,
		ScheduledDate: src.ScheduledDate,
		Status:        model.WorkoutStatusPlanned,
		IsClubTemplate: false,
		CreatedAt:     src.CreatedAt,
	}
	for i, s := range src.Segments {
		w.Segments = append(w.Segments, model.NewSegment(s.Kind, s.Title, s.DistKm, s.Pace, i))
	}
	return w
}
