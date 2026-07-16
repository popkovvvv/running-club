package plan

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/app/http/middleware"
	"github.com/nikpopkov/running-club/api/internal/app/http/response"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type useCase interface {
	ListWeeks(ctx context.Context, coachID uuid.UUID) ([]dto.PlanWeekView, error)
	UpsertWeek(ctx context.Context, coachID uuid.UUID, weekIndex int, req dto.UpsertPlanWeekRequest) (*dto.PlanWeekView, error)
	GetTemplate(ctx context.Context, coachID uuid.UUID, weekIndex int) (*dto.TemplateResponse, error)
	SaveTemplate(ctx context.Context, coachID uuid.UUID, weekIndex int, req dto.SaveTemplateRequest) (*dto.TemplateResponse, error)
	Publish(ctx context.Context, coachID uuid.UUID, weekIndex int) error
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ListWeeks(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	res, err := h.uc.ListWeeks(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) UpsertWeek(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	weekIndex, err := strconv.Atoi(chi.URLParam(r, "weekIndex"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid weekIndex", Code: "bad_request"})
		return
	}
	var req dto.UpsertPlanWeekRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.UpsertWeek(r.Context(), middleware.UserID(r.Context()), weekIndex, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	weekIndex, err := strconv.Atoi(chi.URLParam(r, "weekIndex"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid weekIndex", Code: "bad_request"})
		return
	}
	res, err := h.uc.GetTemplate(r.Context(), middleware.UserID(r.Context()), weekIndex)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) SaveTemplate(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	weekIndex, err := strconv.Atoi(chi.URLParam(r, "weekIndex"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid weekIndex", Code: "bad_request"})
		return
	}
	var req dto.SaveTemplateRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.SaveTemplate(r.Context(), middleware.UserID(r.Context()), weekIndex, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	weekIndex, err := strconv.Atoi(chi.URLParam(r, "weekIndex"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid weekIndex", Code: "bad_request"})
		return
	}
	if err := h.uc.Publish(r.Context(), middleware.UserID(r.Context()), weekIndex); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
