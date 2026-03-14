package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
)

// AdminHandler handles admin endpoints for owner-to-Sleeper mapping CRUD.
type AdminHandler struct {
	svc *service.SleeperService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(svc *service.SleeperService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// RegisterRoutes registers admin endpoints under /api/v1/admin.
func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/admin/mappings", func(r chi.Router) {
		r.Get("/", h.ListMappings)
		r.Post("/", h.CreateMapping)
		r.Put("/{id}", h.UpdateMapping)
		r.Delete("/{id}", h.DeleteMapping)
	})
}

func (h *AdminHandler) ListMappings(w http.ResponseWriter, r *http.Request) {
	mappings, err := h.svc.ListMappings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list mappings")
		return
	}
	writeJSON(w, http.StatusOK, mappings)
}

type createMappingRequest struct {
	OwnerID       int    `json:"owner_id"`
	SleeperUserID string `json:"sleeper_user_id"`
}

func (h *AdminHandler) CreateMapping(w http.ResponseWriter, r *http.Request) {
	var req createMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.OwnerID == 0 || req.SleeperUserID == "" {
		writeError(w, http.StatusBadRequest, "owner_id and sleeper_user_id are required")
		return
	}

	mapping, err := h.svc.CreateMapping(r.Context(), req.OwnerID, req.SleeperUserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create mapping")
		return
	}
	writeJSON(w, http.StatusCreated, mapping)
}

type updateMappingRequest struct {
	SleeperUserID string `json:"sleeper_user_id"`
}

func (h *AdminHandler) UpdateMapping(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid mapping id")
		return
	}

	var req updateMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.SleeperUserID == "" {
		writeError(w, http.StatusBadRequest, "sleeper_user_id is required")
		return
	}

	mapping, err := h.svc.UpdateMapping(r.Context(), id, req.SleeperUserID)
	if err != nil {
		writeError(w, http.StatusNotFound, "mapping not found")
		return
	}
	writeJSON(w, http.StatusOK, mapping)
}

func (h *AdminHandler) DeleteMapping(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid mapping id")
		return
	}

	if err := h.svc.DeleteMapping(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "mapping not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
