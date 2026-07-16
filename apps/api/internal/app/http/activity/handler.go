package activity

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/app/http/middleware"
	"github.com/nikpopkov/running-club/api/internal/app/http/response"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type useCase interface {
	ListActivities(ctx context.Context, userID uuid.UUID) ([]dto.ActivityView, error)
	GetByID(ctx context.Context, actorID uuid.UUID, role model.Role, activityID uuid.UUID) (*dto.ActivityDetailView, error)
	GetStreams(ctx context.Context, actorID uuid.UUID, role model.Role, activityID uuid.UUID) ([]dto.ActivityStreamView, error)
	Progress(ctx context.Context, userID uuid.UUID) (*dto.ProgressResponse, error)
	PRs(ctx context.Context, userID uuid.UUID) ([]dto.PRView, error)
	Races(ctx context.Context, userID uuid.UUID) ([]dto.RaceView, error)
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.ListActivities(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	role := model.Role(middleware.Role(r.Context()))
	res, err := h.uc.GetByID(r.Context(), middleware.UserID(r.Context()), role, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Streams(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	role := model.Role(middleware.Role(r.Context()))
	res, err := h.uc.GetStreams(r.Context(), middleware.UserID(r.Context()), role, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Progress(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.Progress(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) PRs(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.PRs(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Races(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.Races(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
