package dto

type StudentDetailView struct {
	Student          StudentView       `json:"student"`
	WeekKm           string            `json:"weekKm"`
	WeekPlanKm       string            `json:"weekPlanKm"`
	Comp             int               `json:"comp"`
	PlanDays         []PlanDayView     `json:"planDays"`
	RecentActivities []ActivityView    `json:"recentActivities"`
	Summary          PeriodSummaryView `json:"summary"`
	Year             int               `json:"year"`
	Month            int               `json:"month"`
}

type PlanDayView struct {
	WorkoutID     string  `json:"workoutId"`
	DayLabel      string  `json:"dayLabel"`
	Title         string  `json:"title"`
	WorkoutType   string  `json:"workoutType"`
	PlannedKm     float64 `json:"plannedKm"`
	ActualKm      float64 `json:"actualKm"`
	Status        string  `json:"status"`
	ScheduledDate string  `json:"scheduledDate,omitempty"`
	ActivityID    string  `json:"activityId,omitempty"`
	ActivityWhen  string  `json:"activityWhen,omitempty"`
	ActivityPace  string  `json:"activityPace,omitempty"`
}

func NewStudentDetailView(
	student StudentView,
	weekKm, weekPlanKm string,
	comp int,
	planDays []PlanDayView,
	recentActivities []ActivityView,
	summary PeriodSummaryView,
	year, month int,
) StudentDetailView {
	return StudentDetailView{
		Student:          student,
		WeekKm:           weekKm,
		WeekPlanKm:       weekPlanKm,
		Comp:             comp,
		PlanDays:         planDays,
		RecentActivities: recentActivities,
		Summary:          summary,
		Year:             year,
		Month:            month,
	}
}

func NewPlanDayView(
	workoutID, dayLabel, title, workoutType string,
	plannedKm, actualKm float64,
	status, scheduledDate string,
	activityID, activityWhen, activityPace string,
) PlanDayView {
	return PlanDayView{
		WorkoutID:     workoutID,
		DayLabel:      dayLabel,
		Title:         title,
		WorkoutType:   workoutType,
		PlannedKm:     plannedKm,
		ActualKm:      actualKm,
		Status:        status,
		ScheduledDate: scheduledDate,
		ActivityID:    activityID,
		ActivityWhen:  activityWhen,
		ActivityPace:  activityPace,
	}
}
