package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
)

type Handler struct {
	svc *service.HistoricalService
}

func New(svc *service.HistoricalService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/owners", h.ListOwners)
		r.Get("/owners/{id}/stats", h.GetOwnerStats)

		r.Get("/seasons/{year}/standings", h.GetSeasonStandings)
		r.Get("/seasons/{year}/awards", h.GetSeasonAwards)

		r.Get("/belt/current", h.GetCurrentBeltHolder)
		r.Get("/belt/leaderboard", h.GetBeltLeaderboard)
	})
}

func (h *Handler) ListOwners(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	owners, err := h.svc.ListOwners(r.Context(), activeOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list owners")
		return
	}
	writeJSON(w, http.StatusOK, owners)
}

func (h *Handler) GetOwnerStats(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid owner id")
		return
	}

	stats, err := h.svc.GetOwnerStats(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "owner not found")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetSeasonStandings(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year")
		return
	}

	standings, err := h.svc.GetSeasonStandings(r.Context(), year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get standings")
		return
	}
	writeJSON(w, http.StatusOK, standings)
}

func (h *Handler) GetSeasonAwards(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year")
		return
	}

	awards, err := h.svc.GetSeasonAwards(r.Context(), year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get awards")
		return
	}
	writeJSON(w, http.StatusOK, awards)
}

func (h *Handler) GetCurrentBeltHolder(w http.ResponseWriter, r *http.Request) {
	belt, err := h.svc.GetCurrentBeltHolder(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, "no belt history found")
		return
	}
	writeJSON(w, http.StatusOK, belt)
}

func (h *Handler) GetBeltLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := h.svc.GetBeltLeaderboard(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get belt leaderboard")
		return
	}
	writeJSON(w, http.StatusOK, leaderboard)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
