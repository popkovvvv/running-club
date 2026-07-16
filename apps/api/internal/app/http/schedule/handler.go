package schedule

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
	List(ctx context.Context, userID uuid.UUID, role string) ([]dto.AnnounceView, error)
	Publish(ctx context.Context, coachID uuid.UUID, req dto.CreateAnnounceRequest) (*dto.AnnounceView, error)
	Signup(ctx context.Context, athleteID, announceID uuid.UUID) (*dto.AnnounceView, error)
	Unsignup(ctx context.Context, athleteID, announceID uuid.UUID) (*dto.AnnounceView, error)
	Calendar(ctx context.Context, userID uuid.UUID, role string, year, month int) (*dto.CalendarResponse, error)
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.List(r.Context(), middleware.UserID(r.Context()), middleware.Role(r.Context()))
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
	var req dto.CreateAnnounceRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.Publish(r.Context(), middleware.UserID(r.Context()), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	res, err := h.uc.Signup(r.Context(), middleware.UserID(r.Context()), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Unsignup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	res, err := h.uc.Unsignup(r.Context(), middleware.UserID(r.Context()), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Calendar(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	res, err := h.uc.Calendar(r.Context(), middleware.UserID(r.Context()), middleware.Role(r.Context()), year, month)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
