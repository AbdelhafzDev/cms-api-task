package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"cms-api/internal/modules/importer/dto"
	"cms-api/internal/modules/importer/service"
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

func (h *Handler) ListSources(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListSources(r.Context())
	if err != nil {
		h.log.Error("failed to list import sources", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) RunSource(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	pathID := dto.PathSourceID{ID: id}
	if err := validator.Validate(pathID); err != nil {
		httputil.BadRequest(w, "invalid source id")
		return
	}

	resp, err := h.service.RunSource(r.Context(), id)
	if err != nil {
		h.log.Error("failed to run import source", zap.Error(err), zap.Int64("source_id", id))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}
