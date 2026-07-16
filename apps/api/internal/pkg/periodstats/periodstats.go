package periodstats

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func Build(now time.Time, activities []*model.Activity, planWorkouts []*model.Workout) dto.PeriodSummaryView {
	now = now.UTC()
	weekStart := startOfWeek(now)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	weekEnd := weekStart.AddDate(0, 0, 7)
	monthEnd := monthStart.AddDate(0, 1, 0)
	yearEnd := yearStart.AddDate(1, 0, 0)

	return dto.NewPeriodSummaryView(
		statsFor(activities, planWorkouts, weekStart, weekEnd),
		statsFor(activities, planWorkouts, monthStart, monthEnd),
		statsFor(activities, planWorkouts, yearStart, yearEnd),
	)
}

func statsFor(activities []*model.Activity, planWorkouts []*model.Workout, from, to time.Time) dto.PeriodStatsView {
	var km float64
	var count int
	var paceSum float64
	var paceN int

	for _, a := range activities {
		at := activityTime(a)
		if at.Before(from) || !at.Before(to) {
			continue
		}
		km += a.DistKm
		count++
		if sec, ok := paceToSeconds(a.Pace); ok {
			paceSum += float64(sec)
			paceN++
		}
	}

	var planned, completed int
	for _, w := range planWorkouts {
		at := workoutTime(w)
		if at.Before(from) || !at.Before(to) {
			continue
		}
		planned++
		if w.Status == model.WorkoutStatusCompleted {
			completed++
		}
	}

	planCompletion := "—"
	if planned > 0 {
		planCompletion = fmt.Sprintf("%d%%", int(math.Round(100*float64(completed)/float64(planned))))
	}

	avgPace := "—"
	if paceN > 0 {
		avgPace = secondsToPace(int(math.Round(paceSum / float64(paceN))))
	}

	return dto.NewPeriodStatsView(
		strconv.FormatFloat(math.Round(km*10)/10, 'f', 1, 64),
		count,
		planCompletion,
		avgPace,
	)
}

func activityTime(a *model.Activity) time.Time {
	if a.StartedAt != nil {
		return a.StartedAt.UTC()
	}
	return a.CreatedAt.UTC()
}

func workoutTime(w *model.Workout) time.Time {
	if w.ScheduledDate != nil {
		return w.ScheduledDate.UTC()
	}
	return w.CreatedAt.UTC()
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func paceToSeconds(pace string) (int, bool) {
	pace = strings.TrimSpace(pace)
	if pace == "" || pace == "—" {
		return 0, false
	}
	parts := strings.Split(pace, ":")
	if len(parts) != 2 {
		return 0, false
	}
	min, err1 := strconv.Atoi(parts[0])
	sec, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || min < 0 || sec < 0 || sec >= 60 {
		return 0, false
	}
	return min*60 + sec, true
}

func secondsToPace(total int) string {
	if total < 0 {
		total = 0
	}
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}
