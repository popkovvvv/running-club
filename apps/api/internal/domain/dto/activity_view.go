package dto

import "github.com/google/uuid"

type ActivityView struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	When       string    `json:"when"`
	Dist       string    `json:"dist"`
	Time       string    `json:"time"`
	Pace       string    `json:"pace"`
	HR         string    `json:"hr"`
	Kudos      int       `json:"kudos"`
	ComN       int       `json:"comN"`
	Route      string    `json:"route"`
	SX         float64   `json:"sx"`
	SY         float64   `json:"sy"`
	EX         float64   `json:"ex"`
	EY         float64   `json:"ey"`
	Source     string    `json:"source,omitempty"`
	SportType  string    `json:"sportType,omitempty"`
	Elevation  float64   `json:"elevation,omitempty"`
	Visibility string    `json:"visibility,omitempty"`
}

func NewActivityView(
	id uuid.UUID,
	title, when, dist, timeLabel, pace, hr string,
	kudos, comN int,
	route string,
	sx, sy, ex, ey float64,
	source, sportType string,
	elevation float64,
	visibility string,
) ActivityView {
	return ActivityView{
		ID:         id,
		Title:      title,
		When:       when,
		Dist:       dist,
		Time:       timeLabel,
		Pace:       pace,
		HR:         hr,
		Kudos:      kudos,
		ComN:       comN,
		Route:      route,
		SX:         sx,
		SY:         sy,
		EX:         ex,
		EY:         ey,
		Source:     source,
		SportType:  sportType,
		Elevation:  elevation,
		Visibility: visibility,
	}
}

type ProgressResponse struct {
	Months     []MonthStatView `json:"months"`
	YearKm     float64         `json:"yearKm"`
	YearTr     int             `json:"yearTr"`
	YearStarts int             `json:"yearStarts"`
}

func NewProgressResponse(yearKm float64, yearTr, yearStarts int, months []MonthStatView) *ProgressResponse {
	return &ProgressResponse{
		Months:     months,
		YearKm:     yearKm,
		YearTr:     yearTr,
		YearStarts: yearStarts,
	}
}

type MonthStatView struct {
	M    string  `json:"m"`
	Km   float64 `json:"km"`
	Tr   int     `json:"tr"`
	Pace string  `json:"pace"`
	Diff string  `json:"diff"`
}

func NewMonthStatView(month string, km float64, tr int, pace, diff string) MonthStatView {
	return MonthStatView{
		M:    month,
		Km:   km,
		Tr:   tr,
		Pace: pace,
		Diff: diff,
	}
}

type AnalyticsResponse struct {
	ClubKm     float64       `json:"clubKm"`
	Attendance int           `json:"attendance"`
	Students   []StudentView `json:"students"`
}

func NewAnalyticsResponse(clubKm float64, attendance int, students []StudentView) *AnalyticsResponse {
	return &AnalyticsResponse{
		ClubKm:     clubKm,
		Attendance: attendance,
		Students:   students,
	}
}

type PRView struct {
	D    string `json:"d"`
	T    string `json:"t"`
	Date string `json:"date"`
	Col  string `json:"col,omitempty"`
}

func NewPRView(distance, timeLabel, date string) PRView {
	return PRView{
		D:    distance,
		T:    timeLabel,
		Date: date,
	}
}

type RaceView struct {
	Days string `json:"days"`
	Name string `json:"name"`
	Date string `json:"date"`
	Dist string `json:"dist"`
	Goal string `json:"goal"`
}

func NewRaceView(days, name, date, dist, goal string) RaceView {
	return RaceView{
		Days: days,
		Name: name,
		Date: date,
		Dist: dist,
		Goal: goal,
	}
}
