package schedule_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		announceRepo   announceRepo
		clubRepo       clubResolver
		membershipRepo membershipRepo
		workoutRepo    workoutRepo
	}

	announceRepo interface {
		Create(ctx context.Context, a *model.Announce) error
		FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.Announce, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Announce, error)
		IncGoing(ctx context.Context, id uuid.UUID, delta int) error
		CreateSignup(ctx context.Context, s *model.AnnounceSignup) error
		DeleteSignup(ctx context.Context, announceID, athleteID uuid.UUID) error
		HasSignup(ctx context.Context, announceID, athleteID uuid.UUID) (bool, error)
		FindGoingAthletes(ctx context.Context, announceID uuid.UUID) ([]*model.User, error)
	}

	clubResolver interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error)
	}

	membershipRepo interface {
		GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
	}

	workoutRepo interface {
		Create(ctx context.Context, w *model.Workout) error
		DeleteByUserAndAnnounce(ctx context.Context, userID, announceID uuid.UUID) error
	}
)

func NewUseCase(
	announceRepo announceRepo,
	clubRepo clubResolver,
	membershipRepo membershipRepo,
	workoutRepo workoutRepo,
) *UseCase {
	return &UseCase{
		announceRepo:   announceRepo,
		clubRepo:       clubRepo,
		membershipRepo: membershipRepo,
		workoutRepo:    workoutRepo,
	}
}

func workoutFromAnnounce(athleteID uuid.UUID, a *model.Announce) *model.Workout {
	weekIndex := 0
	if a.StartsOn != nil {
		weekIndex = weekIndexForDate(*a.StartsOn, time.Now().UTC())
	}
	clubID := a.ClubID
	annID := a.ID
	w := model.NewWorkout(
		athleteID,
		model.WorkoutPlan,
		a.DayLabel,
		"Групповая",
		a.Place,
		0,
		"",
		weekIndex,
	)
	w.ClubID = &clubID
	w.AnnounceID = &annID
	w.WorkoutType = model.WorkoutTypeCross
	w.Description = a.Note
	if a.GroupName != "" {
		if w.Description != "" {
			w.Description = a.GroupName + " · " + w.Description
		} else {
			w.Description = a.GroupName
		}
	}
	w.ScheduledDate = a.StartsOn
	w.Status = model.WorkoutStatusPlanned
	return w
}

func weekIndexForDate(d, now time.Time) int {
	start := startOfWeek(now)
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	days := int(day.Sub(start).Hours() / 24)
	if days < 0 {
		return 0
	}
	wi := days / 7
	if wi > 3 {
		wi = 3
	}
	return wi
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
