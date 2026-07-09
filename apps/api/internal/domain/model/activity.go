package model

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	WhenLabel string
	DistKm    float64
	Duration  string
	Pace      string
	HR        int
	Kudos     int
	Comments  int
	RouteSVG  string
	StartX    float64
	StartY    float64
	EndX      float64
	EndY      float64
	CreatedAt time.Time
}

func NewActivity(
	userID uuid.UUID,
	title, whenLabel string,
	distKm float64,
	duration, pace string,
	hr, kudos, comments int,
	routeSVG string,
	startX, startY, endX, endY float64,
) *Activity {
	return &Activity{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		WhenLabel: whenLabel,
		DistKm:    distKm,
		Duration:  duration,
		Pace:      pace,
		HR:        hr,
		Kudos:     kudos,
		Comments:  comments,
		RouteSVG:  routeSVG,
		StartX:    startX,
		StartY:    startY,
		EndX:      endX,
		EndY:      endY,
		CreatedAt: time.Now().UTC(),
	}
}

type PR struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Distance  string
	Time      string
	DateLabel string
	Pending   bool
}

func NewPR(userID uuid.UUID, distance, timeLabel, dateLabel string) *PR {
	return &PR{
		ID:        uuid.New(),
		UserID:    userID,
		Distance:  distance,
		Time:      timeLabel,
		DateLabel: dateLabel,
	}
}

type Race struct {
	ID        uuid.UUID
	ClubID    *uuid.UUID
	UserID    *uuid.UUID
	Name      string
	DateLabel string
	Dist      string
	Goal      string
	DaysLeft  int
	Finished  bool
	Result    string
}

func NewRace(userID uuid.UUID, name, dateLabel, dist, goal string, daysLeft int) *Race {
	uid := userID
	return &Race{
		ID:        uuid.New(),
		UserID:    &uid,
		Name:      name,
		DateLabel: dateLabel,
		Dist:      dist,
		Goal:      goal,
		DaysLeft:  daysLeft,
	}
}

type MonthStat struct {
	Month string
	Km    float64
	Tr    int
	Pace  string
	Diff  string
}

func NewMonthStat(month string, km float64, tr int, pace, diff string) MonthStat {
	return MonthStat{
		Month: month,
		Km:    km,
		Tr:    tr,
		Pace:  pace,
		Diff:  diff,
	}
}
