package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rmarko/electronics-marketplace/backend/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) error
	GetByID(ctx context.Context, id int) (*model.Product, error)
	Update(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter model.ProductFilter) ([]model.Product, int, error)
	ListCategories(ctx context.Context) ([]model.Category, error)
}

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) ProductRepository {
	return &productRepo{pool: pool}
}

func (r *productRepo) Create(ctx context.Context, p *model.Product) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO products (seller_id, category_id, brand, model, condition, price, description, image_url, specs, stock)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at, updated_at`,
		p.SellerID, p.CategoryID, p.Brand, p.Model, p.Condition,
		p.Price, p.Description, p.ImageURL, p.Specs, p.Stock,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *productRepo) GetByID(ctx context.Context, id int) (*model.Product, error) {
	p := &model.Product{}
	err := r.pool.QueryRow(ctx,
		`SELECT p.id, p.seller_id, p.category_id, p.brand, p.model, p.condition,
		        p.price, p.description, p.image_url, p.specs, p.stock,
		        p.created_at, p.updated_at,
		        c.name AS category_name, u.name AS seller_name
		 FROM products p
		 JOIN categories c ON c.id = p.category_id
		 JOIN users u ON u.id = p.seller_id
		 WHERE p.id = $1`, id,
	).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.Brand, &p.Model, &p.Condition,
		&p.Price, &p.Description, &p.ImageURL, &p.Specs, &p.Stock,
		&p.CreatedAt, &p.UpdatedAt,
		&p.CategoryName, &p.SellerName,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productRepo) Update(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error) {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if req.Brand != nil {
		setClauses = append(setClauses, fmt.Sprintf("brand = $%d", argIdx))
		args = append(args, *req.Brand)
		argIdx++
	}
	if req.Model != nil {
		setClauses = append(setClauses, fmt.Sprintf("model = $%d", argIdx))
		args = append(args, *req.Model)
		argIdx++
	}
	if req.Condition != nil {
		setClauses = append(setClauses, fmt.Sprintf("condition = $%d", argIdx))
		args = append(args, *req.Condition)
		argIdx++
	}
	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, *req.Price)
		argIdx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}
	if req.ImageURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("image_url = $%d", argIdx))
		args = append(args, *req.ImageURL)
		argIdx++
	}
	if req.Specs != nil {
		setClauses = append(setClauses, fmt.Sprintf("specs = $%d", argIdx))
		args = append(args, *req.Specs)
		argIdx++
	}
	if req.Stock != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock = $%d", argIdx))
		args = append(args, *req.Stock)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE products SET %s WHERE id = $%d
		 RETURNING id, seller_id, category_id, brand, model, condition,
		           price, description, image_url, specs, stock, created_at, updated_at`,
		strings.Join(setClauses, ", "), argIdx,
	)

	p := &model.Product{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.Brand, &p.Model, &p.Condition,
		&p.Price, &p.Description, &p.ImageURL, &p.Specs, &p.Stock,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *productRepo) List(ctx context.Context, filter model.ProductFilter) ([]model.Product, int, error) {
	where := []string{}
	args := []any{}
	argIdx := 1

	if filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("p.category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.Brand != nil {
		where = append(where, fmt.Sprintf("LOWER(p.brand) = LOWER($%d)", argIdx))
		args = append(args, *filter.Brand)
		argIdx++
	}
	if filter.Condition != nil {
		where = append(where, fmt.Sprintf("p.condition = $%d", argIdx))
		args = append(args, *filter.Condition)
		argIdx++
	}
	if filter.MinPrice != nil {
		where = append(where, fmt.Sprintf("p.price >= $%d", argIdx))
		args = append(args, *filter.MinPrice)
		argIdx++
	}
	if filter.MaxPrice != nil {
		where = append(where, fmt.Sprintf("p.price <= $%d", argIdx))
		args = append(args, *filter.MaxPrice)
		argIdx++
	}
	if filter.Search != nil && *filter.Search != "" {
		where = append(where, fmt.Sprintf("(LOWER(p.brand) LIKE LOWER($%d) OR LOWER(p.model) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d))", argIdx, argIdx, argIdx))
		args = append(args, "%"+*filter.Search+"%")
		argIdx++
	}
	if filter.SellerID != nil {
		where = append(where, fmt.Sprintf("p.seller_id = $%d", argIdx))
		args = append(args, *filter.SellerID)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM products p %s`, whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := "p.created_at DESC"
	switch filter.Sort {
	case "price_asc":
		orderBy = "p.price ASC"
	case "price_desc":
		orderBy = "p.price DESC"
	case "newest":
		orderBy = "p.created_at DESC"
	case "oldest":
		orderBy = "p.created_at ASC"
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	args = append(args, filter.PerPage, offset)
	dataQuery := fmt.Sprintf(
		`SELECT p.id, p.seller_id, p.category_id, p.brand, p.model, p.condition,
		        p.price, p.description, p.image_url, p.specs, p.stock,
		        p.created_at, p.updated_at,
		        c.name AS category_name, u.name AS seller_name
		 FROM products p
		 JOIN categories c ON c.id = p.category_id
		 JOIN users u ON u.id = p.seller_id
		 %s
		 ORDER BY %s
		 LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIdx, argIdx+1,
	)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Brand, &p.Model, &p.Condition,
			&p.Price, &p.Description, &p.ImageURL, &p.Specs, &p.Stock,
			&p.CreatedAt, &p.UpdatedAt,
			&p.CategoryName, &p.SellerName,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (r *productRepo) ListCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, slug FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []model.Category{}
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
