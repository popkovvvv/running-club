package club

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
	GetClub(ctx context.Context, userID uuid.UUID, role string) (*dto.ClubView, error)
	Join(ctx context.Context, userID uuid.UUID, code string) (*dto.ClubView, error)
	Leave(ctx context.Context, userID uuid.UUID) error
	UpdatePalette(ctx context.Context, coachID uuid.UUID, accent string) (*dto.ClubView, error)
	ListStudents(ctx context.Context, coachID uuid.UUID) ([]dto.StudentView, error)
	RemoveStudent(ctx context.Context, coachID, studentID uuid.UUID) error
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	res, err := h.uc.GetClub(r.Context(), middleware.UserID(r.Context()), middleware.Role(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	var req dto.JoinClubRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.Join(r.Context(), middleware.UserID(r.Context()), req.Code)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Leave(r.Context(), middleware.UserID(r.Context())); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Palette(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	var req dto.PaletteRequest
	if err := response.Decode(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: err.Error(), Code: "bad_request"})
		return
	}
	res, err := h.uc.UpdatePalette(r.Context(), middleware.UserID(r.Context()), req.AccentHex)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Students(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	res, err := h.uc.ListStudents(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) RemoveStudent(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	if err := h.uc.RemoveStudent(r.Context(), middleware.UserID(r.Context()), id); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
