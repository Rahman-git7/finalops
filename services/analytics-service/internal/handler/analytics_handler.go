package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/finalops/analytics-service/internal/service"
	mw "github.com/finalops/analytics-service/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) Calendar(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if v := r.URL.Query().Get("year"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			year = n
		}
	}
	if v := r.URL.Query().Get("month"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			month = n
		}
	}

	days, err := h.svc.GetCalendar(r.Context(), userID, year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, days)
}

func (h *AnalyticsHandler) Progression(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	exerciseID, err := uuid.Parse(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid exercise_id")
		return
	}
	points, err := h.svc.GetProgression(r.Context(), userID, exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) WeeklyStats(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	stats, err := h.svc.GetWeeklyStats(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) MonthlyStats(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if v := r.URL.Query().Get("year"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			year = n
		}
	}
	if v := r.URL.Query().Get("month"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			month = n
		}
	}

	stats, err := h.svc.GetMonthlyStats(r.Context(), userID, year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
