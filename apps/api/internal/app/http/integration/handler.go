package integration

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/app/http/middleware"
	"github.com/nikpopkov/running-club/api/internal/app/http/response"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/usecase/strava_usecase"
)

type useCase interface {
	Status(ctx context.Context, userID uuid.UUID) (*dto.IntegrationView, error)
	ConnectURL(userID uuid.UUID) (string, error)
	CompleteConnect(ctx context.Context, userID uuid.UUID, code string) (*dto.IntegrationView, error)
	Disconnect(ctx context.Context, userID uuid.UUID) error
	HandleWebhook(ctx context.Context, event strava_usecase.WebhookEvent) error
}

type Handler struct {
	uc                 useCase
	webhookVerifyToken string
	webBaseURL         string
}

func NewHandler(uc useCase, webhookVerifyToken, webBaseURL string) *Handler {
	return &Handler{
		uc:                 uc,
		webhookVerifyToken: webhookVerifyToken,
		webBaseURL:         webBaseURL,
	}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.Status(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	url, err := h.uc.ConnectURL(middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Disconnect(r.Context(), middleware.UserID(r.Context())); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.URL.Query().Get("state"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	if _, err := h.uc.CompleteConnect(r.Context(), userID, r.URL.Query().Get("code")); err != nil {
		http.Redirect(w, r, h.webBaseURL+"?strava=error", http.StatusFound)
		return
	}
	http.Redirect(w, r, h.webBaseURL+"?strava=connected", http.StatusFound)
}

func (h *Handler) WebhookVerify(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.verify_token") != h.webhookVerifyToken {
		response.JSON(w, http.StatusForbidden, response.ErrorBody{Error: "forbidden", Code: "forbidden"})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"hub.challenge": r.URL.Query().Get("hub.challenge")})
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	var event strava_usecase.WebhookEvent
	if err := response.Decode(r, &event); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	if err := h.uc.HandleWebhook(r.Context(), event); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
