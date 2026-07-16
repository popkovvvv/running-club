package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nikpopkov/running-club/api/internal/usecase/strava_usecase"
)

const (
	oauthBaseURL = "https://www.strava.com"
	apiBaseURL   = "https://www.strava.com/api/v3"
)

type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

func NewClient(httpClient *http.Client, clientID, clientSecret string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		httpClient:   httpClient,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *Client) ExchangeToken(ctx context.Context, code string) (*strava_usecase.TokenExchange, error) {
	values := url.Values{}
	values.Set("client_id", c.clientID)
	values.Set("client_secret", c.clientSecret)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthBaseURL+"/oauth/token", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
		Scope        string `json:"scope"`
	}
	if err := c.doJSON(req, &resBody); err != nil {
		return nil, fmt.Errorf("doJSON: %w", err)
	}
	scopes := []string{}
	if strings.TrimSpace(resBody.Scope) != "" {
		scopes = strings.Split(resBody.Scope, ",")
	}
	return &strava_usecase.TokenExchange{
		AccessToken:  resBody.AccessToken,
		RefreshToken: resBody.RefreshToken,
		ExpiresAt:    time.Unix(resBody.ExpiresAt, 0).UTC(),
		Scopes:       scopes,
	}, nil
}

func (c *Client) GetAthlete(ctx context.Context, accessToken string) (*strava_usecase.Athlete, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBaseURL+"/athlete", nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	var resBody struct {
		ID int64 `json:"id"`
	}
	if err := c.doJSON(req, &resBody); err != nil {
		return nil, fmt.Errorf("doJSON: %w", err)
	}
	return &strava_usecase.Athlete{ID: resBody.ID}, nil
}

func (c *Client) ListActivities(ctx context.Context, accessToken string, page, perPage int) ([]strava_usecase.ActivitySummary, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiBaseURL+"/athlete/activities?page="+strconv.Itoa(page)+"&per_page="+strconv.Itoa(perPage),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	var resBody []struct {
		ID                 int64     `json:"id"`
		Name               string    `json:"name"`
		SportType          string    `json:"sport_type"`
		Distance           float64   `json:"distance"`
		MovingTime         int       `json:"moving_time"`
		ElapsedTime        int       `json:"elapsed_time"`
		TotalElevationGain float64   `json:"total_elevation_gain"`
		AverageHeartrate   float64   `json:"average_heartrate"`
		MaxHeartrate       float64   `json:"max_heartrate"`
		KudosCount         int       `json:"kudos_count"`
		CommentCount       int       `json:"comment_count"`
		Visibility         string    `json:"visibility"`
		StartDate          time.Time `json:"start_date"`
		Map                struct {
			SummaryPolyline string `json:"summary_polyline"`
		} `json:"map"`
	}
	if err := c.doJSON(req, &resBody); err != nil {
		return nil, fmt.Errorf("doJSON: %w", err)
	}
	out := make([]strava_usecase.ActivitySummary, 0, len(resBody))
	for _, item := range resBody {
		out = append(out, strava_usecase.ActivitySummary{
			ID:                 item.ID,
			Name:               item.Name,
			SportType:          item.SportType,
			Distance:           item.Distance,
			MovingTime:         item.MovingTime,
			ElapsedTime:        item.ElapsedTime,
			TotalElevationGain: item.TotalElevationGain,
			AverageHeartrate:   item.AverageHeartrate,
			MaxHeartrate:       item.MaxHeartrate,
			KudosCount:         item.KudosCount,
			CommentCount:       item.CommentCount,
			Visibility:         item.Visibility,
			StartDate:          item.StartDate,
			Map: strava_usecase.ActivityMap{
				SummaryPolyline: item.Map.SummaryPolyline,
			},
		})
	}
	return out, nil
}

func (c *Client) GetActivityStreams(ctx context.Context, accessToken string, activityID int64) ([]strava_usecase.ActivityStreamPayload, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiBaseURL+"/activities/"+strconv.FormatInt(activityID, 10)+"/streams?keys=latlng,time,heartrate,cadence,altitude&key_by_type=true",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	var resBody map[string]struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := c.doJSON(req, &resBody); err != nil {
		return nil, fmt.Errorf("doJSON: %w", err)
	}
	streams := make([]strava_usecase.ActivityStreamPayload, 0, len(resBody))
	for key, item := range resBody {
		streamType := item.Type
		if streamType == "" {
			streamType = key
		}
		streams = append(streams, strava_usecase.ActivityStreamPayload{
			Type: streamType,
			Data: string(item.Data),
		})
	}
	return streams, nil
}

func (c *Client) doJSON(req *http.Request, dst any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}
	return nil
}
