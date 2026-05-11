package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/finalops/exercise-service/internal/model"
	"github.com/finalops/exercise-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ExerciseHandler struct {
	svc *service.ExerciseService
}

func NewExerciseHandler(svc *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{svc: svc}
}

func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	f := model.ExerciseListFilter{}
	if v := r.URL.Query().Get("category"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 16); err == nil {
			cat := int16(id)
			f.CategoryID = &cat
		}
	}
	if v := r.URL.Query().Get("muscle"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 16); err == nil {
			m := int16(id)
			f.MuscleID = &m
		}
	}
	f.Query = r.URL.Query().Get("q")

	exercises, err := h.svc.List(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list exercises")
		return
	}
	writeJSON(w, http.StatusOK, exercises)
}

func (h *ExerciseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ex, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "exercise not found")
		return
	}
	writeJSON(w, http.StatusOK, ex)
}

func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	ex, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create exercise: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ex)
}

func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.UpdateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ex, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update exercise")
		return
	}
	writeJSON(w, http.StatusOK, ex)
}

func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExerciseHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list categories")
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *ExerciseHandler) ListMuscleGroups(w http.ResponseWriter, r *http.Request) {
	mgs, err := h.svc.ListMuscleGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list muscle groups")
		return
	}
	writeJSON(w, http.StatusOK, mgs)
}

func (h *ExerciseHandler) Health(w http.ResponseWriter, r *http.Request) {
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
