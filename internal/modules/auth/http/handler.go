package http

import (
	"net/http"

	"go.uber.org/zap"

	"cms-api/internal/modules/auth/dto"
	"cms-api/internal/modules/auth/service"
	"cms-api/internal/pkg/httputil"
	"cms-api/internal/pkg/validator"
)

type Handler struct {
	service service.Service
	log     *zap.Logger
}

func NewHandler(service service.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := httputil.DecodeJSON(w, r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	resp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		h.log.Error("login failed", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := httputil.DecodeJSON(w, r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	resp, err := h.service.Refresh(r.Context(), &req)
	if err != nil {
		h.log.Error("token refresh failed", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := httputil.DecodeJSON(w, r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	if err := h.service.Logout(r.Context(), &req); err != nil {
		h.log.Error("logout failed", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}
