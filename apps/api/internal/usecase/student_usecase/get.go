package student_usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type (
	UseCase struct {
		clubRepo       clubRepo
		userRepo       userRepo
		membershipRepo membershipRepo
		activityRepo   activityRepo
		workoutRepo    workoutRepo
		planWeekRepo   planWeekRepo
		announceRepo   announceRepo
	}

	clubRepo interface {
		GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error)
	}

	userRepo interface {
		GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	}

	membershipRepo interface {
		GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error)
	}

	activityRepo interface {
		FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Activity, error)
		SumDistByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (float64, error)
	}

	workoutRepo interface {
		FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error)
		SumPlanDistByUserWeek(ctx context.Context, userID uuid.UUID, weekIndex int) (float64, error)
	}

	planWeekRepo interface {
		GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error)
	}

	announceRepo interface {
		NextLabelForAthlete(ctx context.Context, clubID, athleteID uuid.UUID) (string, error)
	}
)

func NewUseCase(
	clubRepo clubRepo,
	userRepo userRepo,
	membershipRepo membershipRepo,
	activityRepo activityRepo,
	workoutRepo workoutRepo,
	planWeekRepo planWeekRepo,
	announceRepo announceRepo,
) *UseCase {
	return &UseCase{
		clubRepo:       clubRepo,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		activityRepo:   activityRepo,
		workoutRepo:    workoutRepo,
		planWeekRepo:   planWeekRepo,
		announceRepo:   announceRepo,
	}
}

func (u *UseCase) Get(ctx context.Context, coachID, studentID uuid.UUID, weekIndex int) (*dto.StudentDetailView, error) {
	club, err := u.clubRepo.GetByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
	}
	membership, err := u.membershipRepo.GetByUserAndClub(ctx, studentID, club.ID)
	if err != nil {
		return nil, fmt.Errorf("membershipRepo.GetByUserAndClub: %w", err)
	}
	if membership.Status != model.MembershipActive {
		return nil, fmt.Errorf("student not active: %w", model.ErrForbidden)
	}
	usr, err := u.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetByID: %w", err)
	}

	weekStart := startOfWeek(time.Now().UTC())
	weekKm, err := u.activityRepo.SumDistByUserSince(ctx, studentID, weekStart)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.SumDistByUserSince: %w", err)
	}
	planKm, err := u.workoutRepo.SumPlanDistByUserWeek(ctx, studentID, weekIndex)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.SumPlanDistByUserWeek: %w", err)
	}
	if planKm == 0 {
		if pw, err := u.planWeekRepo.GetByClubAndIndex(ctx, club.ID, weekIndex); err == nil {
			planKm, _ = pw.TargetKm()
		}
	}
	comp := 0
	if planKm > 0 {
		comp = int(math.Min(100, math.Round(100*weekKm/planKm)))
	}

	sub, err := u.announceRepo.NextLabelForAthlete(ctx, club.ID, studentID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("announceRepo.NextLabelForAthlete: %w", err)
		}
		sub = "Нет записи"
	}

	planWorkouts, err := u.workoutRepo.FindByUserWeek(ctx, studentID, weekIndex, model.WorkoutPlan)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUserWeek: %w", err)
	}
	activities, err := u.activityRepo.FindByUser(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindByUser: %w", err)
	}

	planDays := make([]dto.PlanDayView, 0, len(planWorkouts))
	for _, w := range planWorkouts {
		actualKm := 0.0
		if w.Status == model.WorkoutStatusCompleted {
			actualKm = w.DistKm
		}
		planDays = append(planDays, dto.NewPlanDayView(
			w.DayLabel, w.Title, string(w.WorkoutType), w.DistKm, actualKm, string(w.Status),
		))
	}

	recent := make([]dto.ActivityView, 0)
	limit := 10
	for i, a := range activities {
		if i >= limit {
			break
		}
		recent = append(recent, dto.NewActivityView(
			a.ID, a.Title, a.WhenLabel, formatKm(a.DistKm), a.Duration, a.Pace, strconv.Itoa(a.HR),
			a.Kudos, a.Comments, a.RouteSVG, a.StartX, a.StartY, a.EndX, a.EndY,
			a.Source, a.SportType, a.ElevationGain, a.Visibility,
		))
	}

	student := dto.NewStudentView(
		usr.ID, initials(usr.Name), usr.Name, sub,
		strconv.FormatFloat(weekKm, 'f', 1, 64), comp,
	)
	return &dto.StudentDetailView{
		Student:          student,
		WeekKm:           strconv.FormatFloat(weekKm, 'f', 1, 64),
		WeekPlanKm:       strconv.FormatFloat(planKm, 'f', 1, 64),
		Comp:             comp,
		PlanDays:         planDays,
		RecentActivities: recent,
	}, nil
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func formatKm(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	var b strings.Builder
	for i, p := range parts {
		if i >= 2 {
			break
		}
		r, _ := utf8.DecodeRuneInString(p)
		b.WriteRune(r)
	}
	return strings.ToUpper(b.String())
}
