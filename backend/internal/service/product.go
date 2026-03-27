package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/model"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/repository"
)

type ProductService struct {
	products repository.ProductRepository
}

func NewProductService(products repository.ProductRepository) *ProductService {
	return &ProductService{products: products}
}

func (s *ProductService) Create(ctx context.Context, sellerID int, req model.CreateProductRequest) (*model.Product, error) {
	if req.Brand == "" || req.Model == "" {
		return nil, model.ErrBadRequest("brand and model are required")
	}
	if req.Price <= 0 {
		return nil, model.ErrBadRequest("price must be positive")
	}
	if req.Stock < 0 {
		return nil, model.ErrBadRequest("stock cannot be negative")
	}
	if req.CategoryID <= 0 {
		return nil, model.ErrBadRequest("valid category is required")
	}
	if !req.Condition.Valid() {
		return nil, model.ErrBadRequest("condition must be one of: like_new, excellent, good, fair")
	}
	if req.Specs == nil {
		req.Specs = json.RawMessage("{}")
	}

	p := &model.Product{
		SellerID:    sellerID,
		CategoryID:  req.CategoryID,
		Brand:       req.Brand,
		Model:       req.Model,
		Condition:   req.Condition,
		Price:       req.Price,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Specs:       req.Specs,
		Stock:       req.Stock,
	}
	if err := s.products.Create(ctx, p); err != nil {
		return nil, model.ErrInternal("failed to create product")
	}
	return p, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*model.Product, error) {
	p, err := s.products.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound("product not found")
		}
		return nil, model.ErrInternal("failed to get product")
	}
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, id, sellerID int, req model.UpdateProductRequest) (*model.Product, error) {
	if req.Condition != nil && !req.Condition.Valid() {
		return nil, model.ErrBadRequest("condition must be one of: like_new, excellent, good, fair")
	}

	tx, err := s.products.BeginTx(ctx)
	if err != nil {
		return nil, model.ErrInternal("failed to start transaction")
	}
	defer tx.Rollback(ctx)

	existing, err := s.products.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound("product not found")
		}
		return nil, model.ErrInternal("failed to get product")
	}
	if existing.SellerID != sellerID {
		return nil, model.ErrForbidden("you can only update your own products")
	}

	p, err := s.products.UpdateTx(ctx, tx, id, req)
	if err != nil {
		return nil, model.ErrInternal("failed to update product")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, model.ErrInternal("failed to commit transaction")
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id, sellerID int) error {
	tx, err := s.products.BeginTx(ctx)
	if err != nil {
		return model.ErrInternal("failed to start transaction")
	}
	defer tx.Rollback(ctx)

	existing, err := s.products.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotFound("product not found")
		}
		return model.ErrInternal("failed to get product")
	}
	if existing.SellerID != sellerID {
		return model.ErrForbidden("you can only delete your own products")
	}

	hasItems, err := s.products.HasOrderItems(ctx, id)
	if err != nil {
		return model.ErrInternal("failed to check order items")
	}
	if hasItems {
		return model.ErrConflict("cannot delete product that has existing orders")
	}

	if err := s.products.DeleteTx(ctx, tx, id); err != nil {
		return model.ErrInternal("failed to delete product")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.ErrInternal("failed to commit transaction")
	}
	return nil
}

func (s *ProductService) List(ctx context.Context, filter model.ProductFilter) (*model.ProductListResponse, error) {
	products, total, err := s.products.List(ctx, filter)
	if err != nil {
		return nil, model.ErrInternal("failed to list products")
	}
	if products == nil {
		products = []model.Product{}
	}
	return &model.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		PerPage:  filter.PerPage,
	}, nil
}

func (s *ProductService) ListCategories(ctx context.Context) ([]model.Category, error) {
	cats, err := s.products.ListCategories(ctx)
	if err != nil {
		return nil, model.ErrInternal("failed to list categories")
	}
	return cats, nil
}
