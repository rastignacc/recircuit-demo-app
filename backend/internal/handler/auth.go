package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rastignacc/recircuit-demo-app/backend/internal/model"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid request body"))
		return
	}

	resp, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		model.WriteError(w, err)
		return
	}

	setTokenCookie(w, resp.Token)
	model.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid request body"))
		return
	}

	resp, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		model.WriteError(w, err)
		return
	}

	setTokenCookie(w, resp.Token)
	model.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	model.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(24 * time.Hour / time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
