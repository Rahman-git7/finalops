package handler

import (
	"encoding/json"
	"net/http"

	"github.com/finalops/auth-service/internal/middleware"
	"github.com/finalops/auth-service/internal/model"
	"github.com/finalops/auth-service/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" || req.Username == "" {
		writeError(w, http.StatusBadRequest, "email, username and password are required")
		return
	}
	tokens, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusConflict, "registration failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tokens, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tokens, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	_ = h.svc.Logout(r.Context(), req.RefreshToken)
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	profile, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch profile")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	var req model.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	profile, err := h.svc.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update profile")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
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
