package dto

import (
	"time"

	"github.com/google/uuid"
)

type ActivityDetailView struct {
	ID                  uuid.UUID  `json:"id"`
	Title               string     `json:"title"`
	When                string     `json:"when"`
	StartedAt           *time.Time `json:"startedAt,omitempty"`
	Dist                string     `json:"dist"`
	Time                string     `json:"time"`
	Pace                string     `json:"pace"`
	HR                  string     `json:"hr"`
	MaxHeartrate        int        `json:"maxHeartrate,omitempty"`
	MovingSeconds       int        `json:"movingSeconds,omitempty"`
	ElapsedSeconds      int        `json:"elapsedSeconds,omitempty"`
	Kudos               int        `json:"kudos"`
	ComN                int        `json:"comN"`
	Route               string     `json:"route"`
	Polyline            string     `json:"polyline,omitempty"`
	SX                  float64    `json:"sx"`
	SY                  float64    `json:"sy"`
	EX                  float64    `json:"ex"`
	EY                  float64    `json:"ey"`
	Source              string     `json:"source,omitempty"`
	SportType           string     `json:"sportType,omitempty"`
	Elevation           float64    `json:"elevation,omitempty"`
	Visibility          string     `json:"visibility,omitempty"`
	ExternalID          string     `json:"externalId,omitempty"`
	LinkedWorkoutID     *uuid.UUID `json:"linkedWorkoutId,omitempty"`
}

type ActivityStreamView struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type UpdateActivityRequest struct {
	Title         *string  `json:"title"`
	When          *string  `json:"when"`
	DistKm        *float64 `json:"distKm"`
	Duration      *string  `json:"duration"`
	Pace          *string  `json:"pace"`
	HR            *int     `json:"hr"`
	ElevationGain *float64 `json:"elevationGain"`
}

func NewActivityDetailView(
	id uuid.UUID,
	title, when string,
	startedAt *time.Time,
	dist, timeLabel, pace, hr string,
	maxHeartrate, movingSeconds, elapsedSeconds int,
	kudos, comN int,
	route, polyline string,
	sx, sy, ex, ey float64,
	source, sportType string,
	elevation float64,
	visibility, externalID string,
	linkedWorkoutID *uuid.UUID,
) ActivityDetailView {
	return ActivityDetailView{
		ID:              id,
		Title:           title,
		When:            when,
		StartedAt:       startedAt,
		Dist:            dist,
		Time:            timeLabel,
		Pace:            pace,
		HR:              hr,
		MaxHeartrate:    maxHeartrate,
		MovingSeconds:   movingSeconds,
		ElapsedSeconds:  elapsedSeconds,
		Kudos:           kudos,
		ComN:            comN,
		Route:           route,
		Polyline:        polyline,
		SX:              sx,
		SY:              sy,
		EX:              ex,
		EY:              ey,
		Source:          source,
		SportType:       sportType,
		Elevation:       elevation,
		Visibility:      visibility,
		ExternalID:      externalID,
		LinkedWorkoutID: linkedWorkoutID,
	}
}

func NewActivityStreamView(streamType, data string) ActivityStreamView {
	return ActivityStreamView{Type: streamType, Data: data}
}
