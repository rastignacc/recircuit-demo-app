package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rmarko/electronics-marketplace/backend/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx pgx.Tx, order *model.Order) error
	CreateOrderItem(ctx context.Context, tx pgx.Tx, item *model.OrderItem) error
	GetProductForUpdate(ctx context.Context, tx pgx.Tx, productID int) (price float64, stock int, err error)
	DecrementStock(ctx context.Context, tx pgx.Tx, productID int, qty int) error
	GetByID(ctx context.Context, id int) (*model.Order, error)
	ListByBuyer(ctx context.Context, buyerID int) ([]model.Order, error)
	ListBySeller(ctx context.Context, sellerID int) ([]model.Order, error)
	GetSellerStats(ctx context.Context, sellerID int) (*model.SellerStats, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type orderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return &orderRepo{pool: pool}
}

func (r *orderRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *orderRepo) CreateOrder(ctx context.Context, tx pgx.Tx, order *model.Order) error {
	return tx.QueryRow(ctx,
		`INSERT INTO orders (buyer_id, status, total)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		order.BuyerID, order.Status, order.Total,
	).Scan(&order.ID, &order.CreatedAt)
}

func (r *orderRepo) CreateOrderItem(ctx context.Context, tx pgx.Tx, item *model.OrderItem) error {
	return tx.QueryRow(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, unit_price)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		item.OrderID, item.ProductID, item.Quantity, item.UnitPrice,
	).Scan(&item.ID)
}

func (r *orderRepo) GetProductForUpdate(ctx context.Context, tx pgx.Tx, productID int) (float64, int, error) {
	var price float64
	var stock int
	err := tx.QueryRow(ctx,
		`SELECT price, stock FROM products WHERE id = $1 FOR UPDATE`,
		productID,
	).Scan(&price, &stock)
	return price, stock, err
}

func (r *orderRepo) DecrementStock(ctx context.Context, tx pgx.Tx, productID int, qty int) error {
	tag, err := tx.Exec(ctx,
		`UPDATE products SET stock = stock - $1, updated_at = NOW()
		 WHERE id = $2 AND stock >= $1`,
		qty, productID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}
	return nil
}

func (r *orderRepo) GetByID(ctx context.Context, id int) (*model.Order, error) {
	o := &model.Order{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, buyer_id, status, total, created_at FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.BuyerID, &o.Status, &o.Total, &o.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price,
		        CONCAT(p.brand, ' ', p.model) AS product_name
		 FROM order_items oi
		 JOIN products p ON p.id = oi.product_id
		 WHERE oi.order_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.ProductName); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}
	return o, nil
}

func (r *orderRepo) ListByBuyer(ctx context.Context, buyerID int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, buyer_id, status, total, created_at
		 FROM orders WHERE buyer_id = $1 ORDER BY created_at DESC`,
		buyerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.BuyerID, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepo) ListBySeller(ctx context.Context, sellerID int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT o.id, o.buyer_id, o.status, o.total, o.created_at
		 FROM orders o
		 JOIN order_items oi ON oi.order_id = o.id
		 JOIN products p ON p.id = oi.product_id
		 WHERE p.seller_id = $1
		 ORDER BY o.created_at DESC`,
		sellerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.BuyerID, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepo) GetSellerStats(ctx context.Context, sellerID int) (*model.SellerStats, error) {
	stats := &model.SellerStats{}

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM products WHERE seller_id = $1`, sellerID,
	).Scan(&stats.TotalListings)
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(oi.quantity), 0), COALESCE(SUM(oi.quantity * oi.unit_price), 0)
		 FROM order_items oi
		 JOIN products p ON p.id = oi.product_id
		 WHERE p.seller_id = $1`,
		sellerID,
	).Scan(&stats.TotalSold, &stats.TotalRevenue)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
