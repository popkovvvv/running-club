package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type ErrorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "internal"
	msg := err.Error()

	switch {
	case errors.Is(err, model.ErrNotFound):
		status, code = http.StatusNotFound, "not_found"
	case errors.Is(err, model.ErrUnauthorized), errors.Is(err, model.ErrInvalidCredentials):
		status, code = http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, model.ErrForbidden):
		status, code = http.StatusForbidden, "forbidden"
	case errors.Is(err, model.ErrConflict), errors.Is(err, model.ErrEmailTaken),
		errors.Is(err, model.ErrAlreadyMember), errors.Is(err, model.ErrAlreadySignedUp):
		status, code = http.StatusConflict, "conflict"
	case errors.Is(err, model.ErrWeakPassword), errors.Is(err, model.ErrInvalidInviteCode),
		errors.Is(err, model.ErrInvalidRole), errors.Is(err, model.ErrNotMember),
		errors.Is(err, model.ErrNotSignedUp), errors.Is(err, model.ErrBadRequest):
		status, code = http.StatusBadRequest, "bad_request"
	}

	JSON(w, status, ErrorBody{Error: msg, Code: code})
}

func Decode(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
