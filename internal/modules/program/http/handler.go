package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"cms-api/internal/modules/program/dto"
	"cms-api/internal/modules/program/service"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProgramRequest
	if err := httputil.DecodeJSON(w, r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	resp, err := h.service.Create(r.Context(), &req)
	if err != nil {
		h.log.Error("failed to create program", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.Created(w, resp)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	pathID := dto.PathID{ID: chi.URLParam(r, "id")}
	if err := validator.Validate(pathID); err != nil {
		httputil.BadRequest(w, "invalid program id")
		return
	}

	var req dto.UpdateProgramRequest
	if err := httputil.DecodeJSON(w, r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	resp, err := h.service.Update(r.Context(), pathID.ID, &req)
	if err != nil {
		h.log.Error("failed to update program", zap.Error(err), zap.String("id", pathID.ID))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	pathID := dto.PathID{ID: chi.URLParam(r, "id")}
	if err := validator.Validate(pathID); err != nil {
		httputil.BadRequest(w, "invalid program id")
		return
	}

	if err := h.service.Delete(r.Context(), pathID.ID); err != nil {
		h.log.Error("failed to delete program", zap.Error(err), zap.String("id", pathID.ID))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	pathID := dto.PathID{ID: chi.URLParam(r, "id")}
	if err := validator.Validate(pathID); err != nil {
		httputil.BadRequest(w, "invalid program id")
		return
	}

	resp, err := h.service.GetByID(r.Context(), pathID.ID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	req := dto.NewListProgramsRequest(cursorStr, limit)

	resp, err := h.service.List(r.Context(), req.Cursor, req.Limit)
	if err != nil {
		h.log.Error("failed to list programs", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}
