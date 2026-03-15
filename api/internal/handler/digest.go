package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
)

// DigestHandler handles the weekly digest endpoint.
type DigestHandler struct {
	svc *service.DigestService
}

// NewDigestHandler creates a new DigestHandler.
func NewDigestHandler(svc *service.DigestService) *DigestHandler {
	return &DigestHandler{svc: svc}
}

// RegisterRoutes registers digest routes on the router.
func (h *DigestHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/digest/weekly", h.GetWeeklyDigest)
	})
}

// GetWeeklyDigest returns the combined weekly summary.
func (h *DigestHandler) GetWeeklyDigest(w http.ResponseWriter, r *http.Request) {
	digest, err := h.svc.GetWeeklyDigest(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get weekly digest")
		return
	}
	writeJSON(w, http.StatusOK, digest)
}
