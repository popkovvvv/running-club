package strava_usecase

import "time"

type TokenExchange struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	Scopes       []string
}

type Athlete struct {
	ID int64
}

type ActivityMap struct {
	SummaryPolyline string
}

type ActivitySummary struct {
	ID                 int64
	Name               string
	SportType          string
	Distance           float64
	MovingTime         int
	ElapsedTime        int
	TotalElevationGain float64
	AverageHeartrate   float64
	MaxHeartrate       float64
	KudosCount         int
	CommentCount       int
	Visibility         string
	StartDate          time.Time
	Map                ActivityMap
}

type ActivityStreamPayload struct {
	Type string
	Data string
}

type WebhookEvent struct {
	ObjectType string
	ObjectID   int64
	AspectType string
	OwnerID    int64
	Updates    map[string]string
}
