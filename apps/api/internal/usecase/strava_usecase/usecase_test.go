//go:build unit

package strava_usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/strava_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/strava_usecase/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStatusDisconnectedWhenIntegrationMissing(t *testing.T) {
	t.Parallel()

	uc := strava_usecase.NewUseCase(nil, nil, nil, nil, strava_usecase.Config{})

	view, err := uc.Status(context.Background(), uuid.New())

	require.NoError(t, err)
	require.Equal(t, "strava", view.Provider)
	require.Equal(t, "disconnected", view.Status)
	require.False(t, view.Connected)
}

func TestStatusConnectedFromStoredIntegration(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	lastSyncedAt := time.Date(2026, time.July, 16, 9, 30, 0, 0, time.UTC)
	expiresAt := time.Date(2026, time.July, 20, 9, 30, 0, 0, time.UTC)
	integration := &model.UserIntegration{
		ID:                uuid.New(),
		UserID:            userID,
		Provider:          model.IntegrationProviderStrava,
		Status:            model.IntegrationStatusActive,
		ExternalAthleteID: "424242",
		Scopes:            []string{"activity:read_all"},
		AccessToken:       "token",
		RefreshToken:      "refresh",
		ExpiresAt:         expiresAt,
		LastSyncedAt:      &lastSyncedAt,
	}

	integrationRepo := mocks.NewIntegrationRepo(t)
	integrationRepo.EXPECT().GetByUserProvider(mock.Anything, userID, model.IntegrationProviderStrava).Return(integration, nil).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, nil, nil, nil, strava_usecase.Config{})

	view, err := uc.Status(context.Background(), userID)

	require.NoError(t, err)
	require.Equal(t, "strava", view.Provider)
	require.Equal(t, "active", view.Status)
	require.True(t, view.Connected)
	require.Equal(t, "424242", view.ExternalAthleteID)
	require.Equal(t, []string{"activity:read_all"}, view.Scopes)
	require.Equal(t, expiresAt, view.ExpiresAt)
	require.NotNil(t, view.LastSyncedAt)
	require.Equal(t, lastSyncedAt, *view.LastSyncedAt)
}

func TestConnectURLIncludesRequiredStravaParams(t *testing.T) {
	t.Parallel()

	cfg := strava_usecase.Config{
		ClientID:    "12345",
		RedirectURL: "http://localhost:8080/api/v1/integrations/strava/callback",
	}
	uc := strava_usecase.NewUseCase(nil, nil, nil, nil, cfg)

	url, err := uc.ConnectURL(uuid.MustParse("7b70f62b-1adb-4696-b071-2480a9fb6d93"))

	require.NoError(t, err)
	require.Contains(t, url, "https://www.strava.com/oauth/authorize")
	require.Contains(t, url, "client_id=12345")
	require.Contains(t, url, "response_type=code")
	require.Contains(t, url, "approval_prompt=auto")
	require.Contains(t, url, "scope=activity%3Aread_all")
	require.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fv1%2Fintegrations%2Fstrava%2Fcallback")
	require.Contains(t, url, "state=7b70f62b-1adb-4696-b071-2480a9fb6d93")
}

func TestCompleteConnectStoresIntegrationAndImportsActivities(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	expiresAt := time.Date(2026, time.July, 20, 9, 30, 0, 0, time.UTC)
	integrationRepo := mocks.NewIntegrationRepo(t)
	activityRepo := mocks.NewActivityRepo(t)
	streamRepo := mocks.NewActivityStreamRepo(t)
	stravaClient := mocks.NewStravaClient(t)

	integrationRepo.EXPECT().GetByUserProvider(mock.Anything, userID, model.IntegrationProviderStrava).Return(nil, model.ErrNotFound).Once()
	integrationRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Twice()
	activityRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Once()
	streamRepo.EXPECT().ReplaceByActivityID(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	stravaClient.EXPECT().ExchangeToken(mock.Anything, "code-123").Return(&strava_usecase.TokenExchange{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    expiresAt,
		Scopes:       []string{"activity:read_all"},
	}, nil).Once()
	stravaClient.EXPECT().GetAthlete(mock.Anything, "access-token").Return(&strava_usecase.Athlete{
		ID: 424242,
	}, nil).Once()
	stravaClient.EXPECT().ListActivities(mock.Anything, "access-token", 1, 30).Return([]strava_usecase.ActivitySummary{
		{
			ID:                 9001,
			Name:               "Morning Run",
			SportType:          "Run",
			Distance:           10000,
			MovingTime:         3000,
			ElapsedTime:        3200,
			TotalElevationGain: 80,
			AverageHeartrate:   151,
			MaxHeartrate:       173,
			KudosCount:         4,
			CommentCount:       1,
			Visibility:         "everyone",
			StartDate:          time.Date(2026, time.July, 16, 6, 30, 0, 0, time.UTC),
			Map: strava_usecase.ActivityMap{
				SummaryPolyline: "encoded",
			},
		},
	}, nil).Once()
	stravaClient.EXPECT().GetActivityStreams(mock.Anything, "access-token", int64(9001)).Return([]strava_usecase.ActivityStreamPayload{
		{Type: "latlng", Data: "[[55.1,37.1],[55.2,37.2]]"},
		{Type: "time", Data: "[0,60]"},
	}, nil).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, activityRepo, streamRepo, stravaClient, strava_usecase.Config{
		ClientID:    "12345",
		RedirectURL: "http://localhost:8080/api/v1/integrations/strava/callback",
	})

	view, err := uc.CompleteConnect(context.Background(), userID, "code-123")

	require.NoError(t, err)
	require.True(t, view.Connected)
	require.Equal(t, "active", view.Status)
	require.Equal(t, "424242", view.ExternalAthleteID)
	require.Equal(t, []string{"activity:read_all"}, view.Scopes)
}

func TestDisconnectMarksIntegrationDisconnected(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	integration := &model.UserIntegration{
		ID:       uuid.New(),
		UserID:   userID,
		Provider: model.IntegrationProviderStrava,
		Status:   model.IntegrationStatusActive,
	}
	integrationRepo := mocks.NewIntegrationRepo(t)
	integrationRepo.EXPECT().GetByUserProvider(mock.Anything, userID, model.IntegrationProviderStrava).Return(integration, nil).Once()
	integrationRepo.EXPECT().Upsert(mock.Anything, integration).Return(nil).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, nil, nil, nil, strava_usecase.Config{})

	err := uc.Disconnect(context.Background(), userID)

	require.NoError(t, err)
	require.Equal(t, model.IntegrationStatusDisconnected, integration.Status)
	require.Empty(t, integration.AccessToken)
	require.Empty(t, integration.RefreshToken)
}

func TestHandleWebhookDeleteRemovesActivity(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	integration := &model.UserIntegration{
		ID:                uuid.New(),
		UserID:            userID,
		Provider:          model.IntegrationProviderStrava,
		Status:            model.IntegrationStatusActive,
		ExternalAthleteID: "424242",
		AccessToken:       "access-token",
	}
	integrationRepo := mocks.NewIntegrationRepo(t)
	activityRepo := mocks.NewActivityRepo(t)

	integrationRepo.EXPECT().GetByProviderExternalAthleteID(mock.Anything, model.IntegrationProviderStrava, "424242").Return(integration, nil).Once()
	activityRepo.EXPECT().DeleteByUserSourceExternalID(mock.Anything, userID, "strava", "9001").Return(nil).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, activityRepo, nil, nil, strava_usecase.Config{})

	err := uc.HandleWebhook(context.Background(), strava_usecase.WebhookEvent{
		ObjectType: "activity",
		ObjectID:   9001,
		AspectType: "delete",
		OwnerID:    424242,
	})

	require.NoError(t, err)
}

func TestHandleWebhookDeauthorizeRevokesIntegration(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	integration := &model.UserIntegration{
		ID:                uuid.New(),
		UserID:            userID,
		Provider:          model.IntegrationProviderStrava,
		Status:            model.IntegrationStatusActive,
		ExternalAthleteID: "424242",
		AccessToken:       "access-token",
		RefreshToken:      "refresh-token",
	}
	integrationRepo := mocks.NewIntegrationRepo(t)
	integrationRepo.EXPECT().GetByProviderExternalAthleteID(mock.Anything, model.IntegrationProviderStrava, "424242").Return(integration, nil).Once()
	integrationRepo.EXPECT().Upsert(mock.Anything, integration).Return(nil).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, nil, nil, nil, strava_usecase.Config{})

	err := uc.HandleWebhook(context.Background(), strava_usecase.WebhookEvent{
		ObjectType: "athlete",
		OwnerID:    424242,
		Updates: map[string]string{
			"authorized": "false",
		},
	})

	require.NoError(t, err)
	require.Equal(t, model.IntegrationStatusRevoked, integration.Status)
	require.Empty(t, integration.AccessToken)
	require.Empty(t, integration.RefreshToken)
}

func TestStatusReturnsRepoError(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	integrationRepo := mocks.NewIntegrationRepo(t)
	integrationRepo.EXPECT().GetByUserProvider(mock.Anything, userID, model.IntegrationProviderStrava).Return(nil, errors.New("db down")).Once()

	uc := strava_usecase.NewUseCase(integrationRepo, nil, nil, nil, strava_usecase.Config{})

	_, err := uc.Status(context.Background(), userID)

	require.EqualError(t, err, "integrationRepo.GetByUserProvider: db down")
}
