package workout

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/app/http/middleware"
	"github.com/nikpopkov/running-club/api/internal/app/http/response"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

type useCase interface {
	Plan(ctx context.Context, userID uuid.UUID, week int) (*dto.PlanResponse, error)
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateWorkoutRequest) (*dto.WorkoutView, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Plan(w http.ResponseWriter, r *http.Request) {
	week, _ := strconv.Atoi(r.URL.Query().Get("week"))
	res, err := h.uc.Plan(r.Context(), middleware.UserID(r.Context()), week)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWorkoutRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.Create(r.Context(), middleware.UserID(r.Context()), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
