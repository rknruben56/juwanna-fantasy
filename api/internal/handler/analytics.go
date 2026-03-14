package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/owners/{id}/vs/{opponent_id}", h.GetHeadToHead)
		r.Get("/owners/{id}/rivalries", h.GetAllRivalries)
		r.Get("/rankings/power", h.GetPowerRankings)
		r.Get("/projections/awards", h.GetAwardProjections)
	})
}

func (h *AnalyticsHandler) GetHeadToHead(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid owner id")
		return
	}

	opponentID, err := strconv.Atoi(chi.URLParam(r, "opponent_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid opponent id")
		return
	}

	if ownerID == opponentID {
		writeError(w, http.StatusBadRequest, "owner and opponent must be different")
		return
	}

	rivalry, err := h.svc.GetHeadToHead(r.Context(), ownerID, opponentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "head-to-head record not found")
		return
	}
	writeJSON(w, http.StatusOK, rivalry)
}

func (h *AnalyticsHandler) GetAllRivalries(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid owner id")
		return
	}

	rivalries, err := h.svc.GetAllRivalries(r.Context(), ownerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get rivalries")
		return
	}
	writeJSON(w, http.StatusOK, rivalries)
}

func (h *AnalyticsHandler) GetPowerRankings(w http.ResponseWriter, r *http.Request) {
	rankings, err := h.svc.GetPowerRankings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get power rankings")
		return
	}
	writeJSON(w, http.StatusOK, rankings)
}

func (h *AnalyticsHandler) GetAwardProjections(w http.ResponseWriter, r *http.Request) {
	projections, err := h.svc.GetAwardProjections(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get award projections")
		return
	}
	writeJSON(w, http.StatusOK, projections)
}
