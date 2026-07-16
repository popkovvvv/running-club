package dto

type PeriodStatsView struct {
	Km             string `json:"km"`
	Workouts       int    `json:"workouts"`
	PlanCompletion string `json:"planCompletion"`
	AvgPace        string `json:"avgPace"`
}

type PeriodSummaryView struct {
	Week  PeriodStatsView `json:"week"`
	Month PeriodStatsView `json:"month"`
	Year  PeriodStatsView `json:"year"`
}

func NewPeriodStatsView(km string, workouts int, planCompletion, avgPace string) PeriodStatsView {
	return PeriodStatsView{
		Km:             km,
		Workouts:       workouts,
		PlanCompletion: planCompletion,
		AvgPace:        avgPace,
	}
}

func NewPeriodSummaryView(week, month, year PeriodStatsView) PeriodSummaryView {
	return PeriodSummaryView{Week: week, Month: month, Year: year}
}
