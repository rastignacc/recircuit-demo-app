package model

import (
	"encoding/json"
	"time"
)

type Condition string

const (
	ConditionLikeNew   Condition = "like_new"
	ConditionExcellent Condition = "excellent"
	ConditionGood      Condition = "good"
	ConditionFair      Condition = "fair"
)

type Product struct {
	ID          int              `json:"id"`
	SellerID    int              `json:"seller_id"`
	CategoryID  int              `json:"category_id"`
	Brand       string           `json:"brand"`
	Model       string           `json:"model"`
	Condition   Condition        `json:"condition"`
	Price       float64          `json:"price"`
	Description string           `json:"description"`
	ImageURL    string           `json:"image_url"`
	Specs       json.RawMessage  `json:"specs"`
	Stock       int              `json:"stock"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Joined fields
	CategoryName string `json:"category_name,omitempty"`
	SellerName   string `json:"seller_name,omitempty"`
}

type CreateProductRequest struct {
	CategoryID  int             `json:"category_id"`
	Brand       string          `json:"brand"`
	Model       string          `json:"model"`
	Condition   Condition       `json:"condition"`
	Price       float64         `json:"price"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Specs       json.RawMessage `json:"specs"`
	Stock       int             `json:"stock"`
}

type UpdateProductRequest struct {
	Brand       *string          `json:"brand,omitempty"`
	Model       *string          `json:"model,omitempty"`
	Condition   *Condition       `json:"condition,omitempty"`
	Price       *float64         `json:"price,omitempty"`
	Description *string          `json:"description,omitempty"`
	ImageURL    *string          `json:"image_url,omitempty"`
	Specs       *json.RawMessage `json:"specs,omitempty"`
	Stock       *int             `json:"stock,omitempty"`
}

type ProductFilter struct {
	CategoryID *int
	Brand      *string
	Condition  *Condition
	MinPrice   *float64
	MaxPrice   *float64
	Search     *string
	SellerID   *int
	Page       int
	PerPage    int
	Sort       string
}

type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}
