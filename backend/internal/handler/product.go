package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rastignacc/electronics-marketplace/backend/internal/middleware"
	"github.com/rastignacc/electronics-marketplace/backend/internal/model"
	"github.com/rastignacc/electronics-marketplace/backend/internal/service"
)

type ProductHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := model.ProductFilter{
		Page:    intQuery(r, "page", 1),
		PerPage: intQuery(r, "per_page", 20),
		Sort:    r.URL.Query().Get("sort"),
	}

	if v := r.URL.Query().Get("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.CategoryID = &id
		}
	}
	if v := r.URL.Query().Get("brand"); v != "" {
		filter.Brand = &v
	}
	if v := r.URL.Query().Get("condition"); v != "" {
		c := model.Condition(v)
		filter.Condition = &c
	}
	if v := r.URL.Query().Get("min_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MinPrice = &f
		}
	}
	if v := r.URL.Query().Get("max_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MaxPrice = &f
		}
	}
	if v := r.URL.Query().Get("search"); v != "" {
		filter.Search = &v
	}

	resp, err := h.productSvc.List(r.Context(), filter)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, resp)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid product id"))
		return
	}

	p, err := h.productSvc.GetByID(r.Context(), id)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	var req model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid request body"))
		return
	}

	p, err := h.productSvc.Create(r.Context(), user.UserID, req)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid product id"))
		return
	}

	var req model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid request body"))
		return
	}

	p, err := h.productSvc.Update(r.Context(), id, user.UserID, req)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		model.WriteError(w, model.ErrBadRequest("invalid product id"))
		return
	}

	if err := h.productSvc.Delete(r.Context(), id, user.UserID); err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.productSvc.ListCategories(r.Context())
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, cats)
}

func (h *ProductHandler) ListSellerProducts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	filter := model.ProductFilter{
		SellerID: &user.UserID,
		Page:     intQuery(r, "page", 1),
		PerPage:  intQuery(r, "per_page", 50),
		Sort:     r.URL.Query().Get("sort"),
	}

	resp, err := h.productSvc.List(r.Context(), filter)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	model.WriteJSON(w, http.StatusOK, resp)
}

func intQuery(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil || i < 1 {
		return defaultVal
	}
	return i
}
