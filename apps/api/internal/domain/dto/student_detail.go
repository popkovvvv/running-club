package dto

type StudentDetailView struct {
	Student          StudentView    `json:"student"`
	WeekKm           string         `json:"weekKm"`
	WeekPlanKm       string         `json:"weekPlanKm"`
	Comp             int            `json:"comp"`
	PlanDays         []PlanDayView  `json:"planDays"`
	RecentActivities []ActivityView `json:"recentActivities"`
}

type PlanDayView struct {
	DayLabel  string  `json:"dayLabel"`
	Title     string  `json:"title"`
	WorkoutType string `json:"workoutType"`
	PlannedKm float64 `json:"plannedKm"`
	ActualKm  float64 `json:"actualKm"`
	Status    string  `json:"status"`
}

func NewStudentDetailView(
	student StudentView,
	weekKm, weekPlanKm string,
	comp int,
	planDays []PlanDayView,
	recentActivities []ActivityView,
) StudentDetailView {
	return StudentDetailView{
		Student:          student,
		WeekKm:           weekKm,
		WeekPlanKm:       weekPlanKm,
		Comp:             comp,
		PlanDays:         planDays,
		RecentActivities: recentActivities,
	}
}

func NewPlanDayView(dayLabel, title, workoutType string, plannedKm, actualKm float64, status string) PlanDayView {
	return PlanDayView{
		DayLabel:    dayLabel,
		Title:       title,
		WorkoutType: workoutType,
		PlannedKm:   plannedKm,
		ActualKm:    actualKm,
		Status:      status,
	}
}
