package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/whokashmiri/rel_ayforge/internal/service"
)

type TenantHandler struct {
	service *service.TenantService
}

type CreateTenantRequest struct {
	Name string `json:"name"`
}

func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{service: tenantService}
}

func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request CreateTenantRequest

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	var extra any

	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return
	}

	tenant, err := h.service.Register(
		r.Context(),
		request.Name,
	)

	if errors.Is(err, service.ErrInvalidName) {
		writeError(w, http.StatusBadRequest, "tenant name must contain 1 to 150 characters")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set(
		"Location",
		"/api/v1/tenants/"+tenant.ID,
	)

	writeJson(w, http.StatusCreated, tenant)
}

func (h *TenantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tenant, err := h.service.GetByID(r.Context(), id)

	switch {
	case errors.Is(err, service.ErrInvalidID):
		writeError(w, http.StatusBadRequest, "invalid tenant ID")
		return

	case errors.Is(err, service.ErrTenantNotFound):
		writeError(w, http.StatusNotFound, "tenant not found")
		return

	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJson(w, http.StatusOK, tenant)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJson(w, statusCode, map[string]string{
		"error": message,
	})
}
