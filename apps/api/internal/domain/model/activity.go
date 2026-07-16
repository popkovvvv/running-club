package model

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Source           string
	ExternalID       string
	SportType        string
	Title            string
	WhenLabel        string
	DistKm           float64
	DistanceMeters   float64
	Duration         string
	Pace             string
	MovingSeconds    int
	ElapsedSeconds   int
	HR               int
	AverageHeartrate int
	MaxHeartrate     int
	ElevationGain    float64
	Kudos            int
	Comments         int
	Visibility       string
	Polyline         string
	RouteSVG         string
	StartX           float64
	StartY           float64
	EndX             float64
	EndY             float64
	StartedAt        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
	now := time.Now().UTC()
	return &Activity{
		ID:             uuid.New(),
		UserID:         userID,
		Title:          title,
		WhenLabel:      whenLabel,
		DistKm:         distKm,
		DistanceMeters: distKm * 1000,
		Duration:       duration,
		Pace:           pace,
		HR:             hr,
		Kudos:          kudos,
		Comments:       comments,
		RouteSVG:       routeSVG,
		StartX:         startX,
		StartY:         startY,
		EndX:           endX,
		EndY:           endY,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

type ActivityStream struct {
	ID         uuid.UUID
	ActivityID uuid.UUID
	Type       string
	DataJSON   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewActivityStream(activityID uuid.UUID, streamType, dataJSON string) *ActivityStream {
	now := time.Now().UTC()
	return &ActivityStream{
		ID:         uuid.New(),
		ActivityID: activityID,
		Type:       streamType,
		DataJSON:   dataJSON,
		CreatedAt:  now,
		UpdatedAt:  now,
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
