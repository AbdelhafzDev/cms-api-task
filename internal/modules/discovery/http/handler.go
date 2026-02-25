package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"cms-api/internal/modules/discovery/dto"
	"cms-api/internal/modules/discovery/service"
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

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	programType := r.URL.Query().Get("type")
	category := r.URL.Query().Get("category")
	language := r.URL.Query().Get("language")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	req := dto.NewSearchRequest(q, programType, category, language, page, perPage)

	if err := validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	resp, err := h.service.Search(r.Context(), &req)
	if err != nil {
		h.log.Error("failed to search programs", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	req := dto.NewListRequest(cursorStr, limit)

	resp, err := h.service.List(r.Context(), req.Cursor, req.Limit)
	if err != nil {
		h.log.Error("failed to list programs", zap.Error(err))
		httputil.HandleError(w, r, err)
		return
	}

	httputil.OK(w, resp)
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
