package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rastignacc/electronics-marketplace/backend/internal/model"
	"github.com/rastignacc/electronics-marketplace/backend/internal/repository"
)

type OrderService struct {
	orders repository.OrderRepository
}

func NewOrderService(orders repository.OrderRepository) *OrderService {
	return &OrderService{orders: orders}
}

func (s *OrderService) PlaceOrder(ctx context.Context, buyerID int, req model.CreateOrderRequest) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, model.ErrBadRequest("order must contain at least one item")
	}
	for _, item := range req.Items {
		if item.Quantity < 1 {
			return nil, model.ErrBadRequest("quantity must be at least 1")
		}
	}

	tx, err := s.orders.BeginTx(ctx)
	if err != nil {
		return nil, model.ErrInternal("failed to start transaction")
	}
	defer tx.Rollback(ctx)

	var totalCents int64
	items := make([]model.OrderItem, len(req.Items))

	for i, reqItem := range req.Items {
		price, stock, err := s.orders.GetProductForUpdate(ctx, tx, reqItem.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, model.ErrNotFound(fmt.Sprintf("product %d not found", reqItem.ProductID))
			}
			return nil, model.ErrInternal("failed to lock product")
		}
		if stock < reqItem.Quantity {
			return nil, model.ErrBadRequest(fmt.Sprintf("insufficient stock for product %d (available: %d)", reqItem.ProductID, stock))
		}

		if err := s.orders.DecrementStock(ctx, tx, reqItem.ProductID, reqItem.Quantity); err != nil {
			return nil, model.ErrInternal("failed to update stock")
		}

		items[i] = model.OrderItem{
			ProductID: reqItem.ProductID,
			Quantity:  reqItem.Quantity,
			UnitPrice: price,
		}
		totalCents += int64(price*100) * int64(reqItem.Quantity)
	}

	total := float64(totalCents) / 100

	order := &model.Order{
		BuyerID: buyerID,
		Status:  model.OrderStatusPending,
		Total:   total,
	}
	if err := s.orders.CreateOrder(ctx, tx, order); err != nil {
		return nil, model.ErrInternal("failed to create order")
	}

	for i := range items {
		items[i].OrderID = order.ID
		if err := s.orders.CreateOrderItem(ctx, tx, &items[i]); err != nil {
			return nil, model.ErrInternal("failed to create order item")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, model.ErrInternal("failed to commit transaction")
	}

	order.Items = items
	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, orderID, userID int, role model.Role) (*model.Order, error) {
	order, err := s.orders.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound("order not found")
		}
		return nil, model.ErrInternal("failed to get order")
	}

	if role == model.RoleBuyer && order.BuyerID != userID {
		return nil, model.ErrForbidden("you can only view your own orders")
	}
	if role == model.RoleSeller {
		has, err := s.orders.HasSellerProduct(ctx, orderID, userID)
		if err != nil {
			return nil, model.ErrInternal("failed to check order ownership")
		}
		if !has {
			return nil, model.ErrForbidden("you can only view orders containing your products")
		}
	}
	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID int, role model.Role, page, perPage int) ([]model.Order, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var orders []model.Order
	var err error

	if role == model.RoleSeller {
		orders, err = s.orders.ListBySeller(ctx, userID, perPage, offset)
	} else {
		orders, err = s.orders.ListByBuyer(ctx, userID, perPage, offset)
	}

	if err != nil {
		return nil, model.ErrInternal("failed to list orders")
	}
	if orders == nil {
		orders = []model.Order{}
	}
	return orders, nil
}

func (s *OrderService) GetSellerStats(ctx context.Context, sellerID int) (*model.SellerStats, error) {
	stats, err := s.orders.GetSellerStats(ctx, sellerID)
	if err != nil {
		return nil, model.ErrInternal("failed to get seller stats")
	}
	return stats, nil
}
