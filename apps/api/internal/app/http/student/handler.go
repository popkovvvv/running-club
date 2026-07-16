package student

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
	Get(ctx context.Context, coachID, studentID uuid.UUID, weekIndex int) (*dto.StudentDetailView, error)
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if middleware.Role(r.Context()) != string(model.RoleCoach) {
		response.Error(w, model.ErrForbidden)
		return
	}
	studentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ErrorBody{Error: "invalid id", Code: "bad_request"})
		return
	}
	weekIndex, _ := strconv.Atoi(r.URL.Query().Get("week"))
	res, err := h.uc.Get(r.Context(), middleware.UserID(r.Context()), studentID, weekIndex)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
