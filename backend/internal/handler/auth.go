package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rmarko/electronics-marketplace/backend/internal/model"
	"github.com/rmarko/electronics-marketplace/backend/internal/service"
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

	model.WriteJSON(w, http.StatusOK, resp)
}
