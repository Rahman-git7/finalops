package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/finalops/workout-service/internal/model"
	"github.com/finalops/workout-service/internal/service"
	mw "github.com/finalops/workout-service/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WorkoutHandler struct {
	svc *service.WorkoutService
}

func NewWorkoutHandler(svc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{svc: svc}
}

// Sessions

func (h *WorkoutHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	var req model.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.SessionDate == "" {
		writeError(w, http.StatusBadRequest, "session_date is required")
		return
	}
	session, err := h.svc.CreateSession(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (h *WorkoutHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}
	sessions, err := h.svc.ListSessions(r.Context(), userID, from, to, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

func (h *WorkoutHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	session, err := h.svc.GetSession(r.Context(), id, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "session not found")
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *WorkoutHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.UpdateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	session, err := h.svc.UpdateSession(r.Context(), id, userID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *WorkoutHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteSession(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Session Exercises

func (h *WorkoutHandler) AddExercise(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	var req model.AddExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	se, err := h.svc.AddExercise(r.Context(), sessionID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, se)
}

func (h *WorkoutHandler) RemoveExercise(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	exerciseID, err := uuid.Parse(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}
	if err := h.svc.RemoveExercise(r.Context(), sessionID, exerciseID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Sets

func (h *WorkoutHandler) AddSet(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	exerciseID, err := uuid.Parse(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}
	var req model.CreateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	set, err := h.svc.AddSet(r.Context(), sessionID, exerciseID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, set)
}

func (h *WorkoutHandler) UpdateSet(w http.ResponseWriter, r *http.Request) {
	setID, err := uuid.Parse(chi.URLParam(r, "set_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid set id")
		return
	}
	var req model.UpdateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	set, err := h.svc.UpdateSet(r.Context(), setID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, set)
}

func (h *WorkoutHandler) DeleteSet(w http.ResponseWriter, r *http.Request) {
	setID, err := uuid.Parse(chi.URLParam(r, "set_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid set id")
		return
	}
	if err := h.svc.DeleteSet(r.Context(), setID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Programs

func (h *WorkoutHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	var req model.CreateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p, err := h.svc.CreateProgram(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *WorkoutHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	programs, err := h.svc.ListPrograms(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, programs)
}

func (h *WorkoutHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.svc.GetProgram(r.Context(), id, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "program not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *WorkoutHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.UpdateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p, err := h.svc.UpdateProgram(r.Context(), id, userID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *WorkoutHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	userID := mw.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteProgram(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) ListProgramTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.ListProgramTypes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, types)
}

// Internal API

func (h *WorkoutHandler) InternalSessionRange(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	items, err := h.svc.GetSessionsInRange(r.Context(), userID, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *WorkoutHandler) InternalExerciseHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	exerciseID, err := uuid.Parse(chi.URLParam(r, "exercise_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid exercise_id")
		return
	}
	history, err := h.svc.GetExerciseHistory(r.Context(), userID, exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func (h *WorkoutHandler) Health(w http.ResponseWriter, r *http.Request) {
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
