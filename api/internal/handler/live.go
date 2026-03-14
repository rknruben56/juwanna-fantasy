package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
)

// LiveHandler handles live Sleeper-powered endpoints.
type LiveHandler struct {
	svc *service.SleeperService
}

// NewLiveHandler creates a new LiveHandler.
func NewLiveHandler(svc *service.SleeperService) *LiveHandler {
	return &LiveHandler{svc: svc}
}

// RegisterRoutes registers live endpoints under /api/v1/live.
func (h *LiveHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/live", func(r chi.Router) {
		r.Get("/matchups", h.GetMatchups)
		r.Get("/standings", h.GetStandings)
		r.Get("/belt", h.GetBeltStatus)
	})
}

func (h *LiveHandler) GetMatchups(w http.ResponseWriter, r *http.Request) {
	matchups, err := h.svc.GetLiveMatchups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get live matchups")
		return
	}
	writeJSON(w, http.StatusOK, matchups)
}

func (h *LiveHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	standings, err := h.svc.GetLiveStandings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get live standings")
		return
	}
	writeJSON(w, http.StatusOK, standings)
}

func (h *LiveHandler) GetBeltStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.svc.GetLiveBeltStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, "belt status unavailable")
		return
	}
	writeJSON(w, http.StatusOK, status)
}
