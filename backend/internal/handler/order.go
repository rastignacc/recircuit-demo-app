package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rmarko/electronics-marketplace/backend/internal/middleware"
	"github.com/rmarko/electronics-marketplace/backend/internal/model"
	"github.com/rmarko/electronics-marketplace/backend/internal/service"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid request body"))
		return
	}

	order, err := h.orderSvc.PlaceOrder(r.Context(), user.UserID, req)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid order id"))
		return
	}

	order, err := h.orderSvc.GetByID(r.Context(), id, user.UserID, user.Role)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	orders, err := h.orderSvc.ListOrders(r.Context(), user.UserID, user.Role)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) SellerStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	stats, err := h.orderSvc.GetSellerStats(r.Context(), user.UserID)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, stats)
}
