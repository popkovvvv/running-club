package strava_usecase

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/polyline"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type integrationRepo interface {
	GetByUserProvider(ctx context.Context, userID uuid.UUID, provider model.IntegrationProvider) (*model.UserIntegration, error)
	GetByProviderExternalAthleteID(ctx context.Context, provider model.IntegrationProvider, externalAthleteID string) (*model.UserIntegration, error)
	Upsert(ctx context.Context, integration *model.UserIntegration) error
}

type activityRepo interface {
	Upsert(ctx context.Context, activity *model.Activity) error
	DeleteByUserSourceExternalID(ctx context.Context, userID uuid.UUID, source, externalID string) error
}

type activityStreamRepo interface {
	ReplaceByActivityID(ctx context.Context, activityID uuid.UUID, streams []*model.ActivityStream) error
}

type stravaClient interface {
	ExchangeToken(ctx context.Context, code string) (*TokenExchange, error)
	GetAthlete(ctx context.Context, accessToken string) (*Athlete, error)
	ListActivities(ctx context.Context, accessToken string, page, perPage int) ([]ActivitySummary, error)
	GetActivityStreams(ctx context.Context, accessToken string, activityID int64) ([]ActivityStreamPayload, error)
}

type UseCase struct {
	integrationRepo    integrationRepo
	activityRepo       activityRepo
	activityStreamRepo activityStreamRepo
	stravaClient       stravaClient
	cfg                Config
}

func NewUseCase(
	integrationRepo integrationRepo,
	activityRepo activityRepo,
	activityStreamRepo activityStreamRepo,
	stravaClient stravaClient,
	cfg Config,
) *UseCase {
	return &UseCase{
		integrationRepo:    integrationRepo,
		activityRepo:       activityRepo,
		activityStreamRepo: activityStreamRepo,
		stravaClient:       stravaClient,
		cfg:                cfg,
	}
}

func (u *UseCase) Status(ctx context.Context, userID uuid.UUID) (*dto.IntegrationView, error) {
	if u.integrationRepo == nil {
		return &dto.IntegrationView{
			Provider:  string(model.IntegrationProviderStrava),
			Status:    string(model.IntegrationStatusDisconnected),
			Connected: false,
		}, nil
	}

	integration, err := u.integrationRepo.GetByUserProvider(ctx, userID, model.IntegrationProviderStrava)
	if err != nil {
		if err == model.ErrNotFound {
			return &dto.IntegrationView{
				Provider:  string(model.IntegrationProviderStrava),
				Status:    string(model.IntegrationStatusDisconnected),
				Connected: false,
			}, nil
		}
		return nil, fmt.Errorf("integrationRepo.GetByUserProvider: %w", err)
	}

	return &dto.IntegrationView{
		Provider:          string(integration.Provider),
		Status:            string(integration.Status),
		Connected:         integration.Status == model.IntegrationStatusActive,
		ExternalAthleteID: integration.ExternalAthleteID,
		Scopes:            integration.Scopes,
		ExpiresAt:         integration.ExpiresAt,
		LastSyncedAt:      integration.LastSyncedAt,
		LastWebhookAt:     integration.LastWebhookAt,
		LastError:         integration.LastError,
	}, nil
}

func (u *UseCase) ConnectURL(userID uuid.UUID) (string, error) {
	if strings.TrimSpace(u.cfg.ClientID) == "" {
		return "", fmt.Errorf("strava client id is empty")
	}
	if strings.TrimSpace(u.cfg.RedirectURL) == "" {
		return "", fmt.Errorf("strava redirect url is empty")
	}

	values := url.Values{}
	values.Set("client_id", u.cfg.ClientID)
	values.Set("response_type", "code")
	values.Set("redirect_uri", u.cfg.RedirectURL)
	values.Set("approval_prompt", "auto")
	values.Set("scope", "activity:read_all")
	values.Set("state", userID.String())

	return "https://www.strava.com/oauth/authorize?" + values.Encode(), nil
}

func (u *UseCase) CompleteConnect(ctx context.Context, userID uuid.UUID, code string) (*dto.IntegrationView, error) {
	if u.stravaClient == nil {
		return nil, fmt.Errorf("strava client is nil")
	}
	tokenExchange, err := u.stravaClient.ExchangeToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("stravaClient.ExchangeToken: %w", err)
	}
	athlete, err := u.stravaClient.GetAthlete(ctx, tokenExchange.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("stravaClient.GetAthlete: %w", err)
	}

	integration, err := u.loadOrCreateIntegration(ctx, userID)
	if err != nil {
		return nil, err
	}
	integration.Status = model.IntegrationStatusActive
	integration.ExternalAthleteID = strconv.FormatInt(athlete.ID, 10)
	integration.AccessToken = tokenExchange.AccessToken
	integration.RefreshToken = tokenExchange.RefreshToken
	integration.ExpiresAt = tokenExchange.ExpiresAt
	integration.Scopes = tokenExchange.Scopes
	integration.LastError = ""
	integration.UpdatedAt = time.Now().UTC()
	if err := u.integrationRepo.Upsert(ctx, integration); err != nil {
		return nil, fmt.Errorf("integrationRepo.Upsert: %w", err)
	}

	activities, err := u.stravaClient.ListActivities(ctx, tokenExchange.AccessToken, 1, 30)
	if err != nil {
		return nil, fmt.Errorf("stravaClient.ListActivities: %w", err)
	}
	for _, summary := range activities {
		if err := u.upsertSummaryActivity(ctx, integration, summary); err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()
	integration.LastSyncedAt = &now
	integration.UpdatedAt = now
	if err := u.integrationRepo.Upsert(ctx, integration); err != nil {
		return nil, fmt.Errorf("integrationRepo.Upsert: %w", err)
	}

	return &dto.IntegrationView{
		Provider:          string(integration.Provider),
		Status:            string(integration.Status),
		Connected:         integration.Status == model.IntegrationStatusActive,
		ExternalAthleteID: integration.ExternalAthleteID,
		Scopes:            integration.Scopes,
		ExpiresAt:         integration.ExpiresAt,
		LastSyncedAt:      integration.LastSyncedAt,
		LastWebhookAt:     integration.LastWebhookAt,
		LastError:         integration.LastError,
	}, nil
}

func (u *UseCase) Disconnect(ctx context.Context, userID uuid.UUID) error {
	integration, err := u.integrationRepo.GetByUserProvider(ctx, userID, model.IntegrationProviderStrava)
	if err != nil {
		if err == model.ErrNotFound {
			return nil
		}
		return fmt.Errorf("integrationRepo.GetByUserProvider: %w", err)
	}
	integration.Status = model.IntegrationStatusDisconnected
	integration.AccessToken = ""
	integration.RefreshToken = ""
	integration.LastError = ""
	integration.UpdatedAt = time.Now().UTC()
	if err := u.integrationRepo.Upsert(ctx, integration); err != nil {
		return fmt.Errorf("integrationRepo.Upsert: %w", err)
	}
	return nil
}

func (u *UseCase) HandleWebhook(ctx context.Context, event WebhookEvent) error {
	integration, err := u.integrationRepo.GetByProviderExternalAthleteID(
		ctx,
		model.IntegrationProviderStrava,
		strconv.FormatInt(event.OwnerID, 10),
	)
	if err != nil {
		if err == model.ErrNotFound {
			return nil
		}
		return fmt.Errorf("integrationRepo.GetByProviderExternalAthleteID: %w", err)
	}

	if event.ObjectType == "athlete" && event.Updates["authorized"] == "false" {
		integration.Status = model.IntegrationStatusRevoked
		integration.AccessToken = ""
		integration.RefreshToken = ""
		integration.UpdatedAt = time.Now().UTC()
		if err := u.integrationRepo.Upsert(ctx, integration); err != nil {
			return fmt.Errorf("integrationRepo.Upsert: %w", err)
		}
		return nil
	}

	if event.ObjectType != "activity" {
		return nil
	}
	if event.AspectType == "delete" {
		if err := u.activityRepo.DeleteByUserSourceExternalID(ctx, integration.UserID, "strava", strconv.FormatInt(event.ObjectID, 10)); err != nil {
			return fmt.Errorf("activityRepo.DeleteByUserSourceExternalID: %w", err)
		}
		return nil
	}
	if u.stravaClient == nil {
		return fmt.Errorf("strava client is nil")
	}
	summaries, err := u.stravaClient.ListActivities(ctx, integration.AccessToken, 1, 30)
	if err != nil {
		return fmt.Errorf("stravaClient.ListActivities: %w", err)
	}
	for _, summary := range summaries {
		if summary.ID == event.ObjectID {
			return u.upsertSummaryActivity(ctx, integration, summary)
		}
	}
	return nil
}

func (u *UseCase) loadOrCreateIntegration(ctx context.Context, userID uuid.UUID) (*model.UserIntegration, error) {
	if u.integrationRepo == nil {
		return nil, fmt.Errorf("integration repo is nil")
	}
	integration, err := u.integrationRepo.GetByUserProvider(ctx, userID, model.IntegrationProviderStrava)
	if err == nil {
		return integration, nil
	}
	if err != model.ErrNotFound {
		return nil, fmt.Errorf("integrationRepo.GetByUserProvider: %w", err)
	}
	return model.NewUserIntegration(userID, model.IntegrationProviderStrava), nil
}

func (u *UseCase) upsertSummaryActivity(ctx context.Context, integration *model.UserIntegration, summary ActivitySummary) error {
	activity := mapSummaryToActivity(integration.UserID, summary)
	if err := u.activityRepo.Upsert(ctx, activity); err != nil {
		return fmt.Errorf("activityRepo.Upsert: %w", err)
	}
	if u.activityStreamRepo == nil || u.stravaClient == nil {
		return nil
	}
	streamPayloads, err := u.stravaClient.GetActivityStreams(ctx, integration.AccessToken, summary.ID)
	if err != nil {
		return fmt.Errorf("stravaClient.GetActivityStreams: %w", err)
	}
	streams := make([]*model.ActivityStream, 0, len(streamPayloads))
	for _, payload := range streamPayloads {
		streams = append(streams, model.NewActivityStream(activity.ID, payload.Type, payload.Data))
	}
	if err := u.activityStreamRepo.ReplaceByActivityID(ctx, activity.ID, streams); err != nil {
		return fmt.Errorf("activityStreamRepo.ReplaceByActivityID: %w", err)
	}
	return nil
}

func mapSummaryToActivity(userID uuid.UUID, summary ActivitySummary) *model.Activity {
	startedAt := summary.StartDate.UTC()
	return &model.Activity{
		ID:               uuid.New(),
		UserID:           userID,
		Source:           "strava",
		ExternalID:       strconv.FormatInt(summary.ID, 10),
		SportType:        summary.SportType,
		Title:            summary.Name,
		WhenLabel:        startedAt.Format("02.01.2006 15:04"),
		DistKm:           summary.Distance / 1000,
		DistanceMeters:   summary.Distance,
		Duration:         formatDuration(summary.MovingTime),
		Pace:             formatPace(summary.Distance, summary.MovingTime),
		MovingSeconds:    summary.MovingTime,
		ElapsedSeconds:   summary.ElapsedTime,
		HR:               int(math.Round(summary.AverageHeartrate)),
		AverageHeartrate: int(math.Round(summary.AverageHeartrate)),
		MaxHeartrate:     int(math.Round(summary.MaxHeartrate)),
		ElevationGain:    summary.TotalElevationGain,
		Kudos:            summary.KudosCount,
		Comments:         summary.CommentCount,
		Visibility:       summary.Visibility,
		Polyline:         summary.Map.SummaryPolyline,
		RouteSVG:         routeSVG(summary.Map.SummaryPolyline),
		StartX:           routeStartX(summary.Map.SummaryPolyline),
		StartY:           routeStartY(summary.Map.SummaryPolyline),
		EndX:             routeEndX(summary.Map.SummaryPolyline),
		EndY:             routeEndY(summary.Map.SummaryPolyline),
		StartedAt:        &startedAt,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	rest := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, rest)
}

func formatPace(distanceMeters float64, movingSeconds int) string {
	if distanceMeters <= 0 || movingSeconds <= 0 {
		return ""
	}
	secondsPerKm := int(math.Round(float64(movingSeconds) / (distanceMeters / 1000)))
	return fmt.Sprintf("%d:%02d", secondsPerKm/60, secondsPerKm%60)
}

func routeSVG(encoded string) string {
	svg := polyline.ToSVG(encoded, 300, 140)
	return svg.Path
}

func routeStartX(encoded string) float64 {
	return polyline.ToSVG(encoded, 300, 140).SX
}

func routeStartY(encoded string) float64 {
	return polyline.ToSVG(encoded, 300, 140).SY
}

func routeEndX(encoded string) float64 {
	return polyline.ToSVG(encoded, 300, 140).EX
}

func routeEndY(encoded string) float64 {
	return polyline.ToSVG(encoded, 300, 140).EY
}
