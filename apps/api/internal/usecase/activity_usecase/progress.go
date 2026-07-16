package activity_usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/periodstats"
)

var monthNamesRU = [...]string{
	"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
	"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
}

func (u *UseCase) Progress(ctx context.Context, userID uuid.UUID) (*dto.ProgressResponse, error) {
	activities, err := u.activityRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindByUser: %w", err)
	}

	year := time.Now().UTC().Year()
	type bucket struct {
		km   float64
		tr   int
		pace string
	}
	byMonth := map[int]*bucket{}

	for _, a := range activities {
		at := a.CreatedAt
		if a.StartedAt != nil {
			at = *a.StartedAt
		}
		if at.Year() != year {
			continue
		}
		m := int(at.Month())
		b := byMonth[m]
		if b == nil {
			b = &bucket{}
			byMonth[m] = b
		}
		b.km += a.DistKm
		b.tr++
		if b.pace == "" && a.Pace != "" {
			b.pace = a.Pace
		}
	}

	orphans, err := u.workoutRepo.FindCompletedWithoutActivity(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindCompletedWithoutActivity: %w", err)
	}
	for _, w := range orphans {
		at := w.CreatedAt
		if w.ScheduledDate != nil {
			at = *w.ScheduledDate
		}
		if at.Year() != year {
			continue
		}
		m := int(at.Month())
		b := byMonth[m]
		if b == nil {
			b = &bucket{}
			byMonth[m] = b
		}
		b.km += w.DistKm
		b.tr++
	}

	months := make([]int, 0, len(byMonth))
	for m := range byMonth {
		months = append(months, m)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(months)))

	var yearKm float64
	var yearTr int
	views := make([]dto.MonthStatView, 0, len(months))
	for _, m := range months {
		b := byMonth[m]
		yearKm += b.km
		yearTr += b.tr
		pace := b.pace
		if pace == "" {
			pace = "—"
		}
		views = append(views, dto.NewMonthStatView(
			monthNamesRU[m],
			math.Round(b.km*10)/10,
			b.tr,
			pace,
			"—",
		))
	}

	races, err := u.activityRepo.FindRaces(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("activityRepo.FindRaces: %w", err)
	}

	allWorkouts, err := u.workoutRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("workoutRepo.FindByUser: %w", err)
	}
	planWorkouts := make([]*model.Workout, 0, len(allWorkouts))
	for _, w := range allWorkouts {
		if w.Kind == model.WorkoutPlan {
			planWorkouts = append(planWorkouts, w)
		}
	}
	summary := periodstats.Build(time.Now().UTC(), activities, planWorkouts)

	return dto.NewProgressResponse(math.Round(yearKm*10)/10, yearTr, len(races), views, summary), nil
}
