package entity

import (
	"strconv"
	"strings"
	"time"
)

type CreateProductRequest struct {
	UserId string `validate:"required,uuid"`

	ShopId      string  `json:"shop_id" validate:"required,uuid"`
	CategoryId  string  `json:"category_id" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,max=255,min=3"`
	Description *string `json:"description" validate:"omitempty,max=255,min=3"`
	ImageUrl    *string `json:"image_url" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,numeric"`
	Stock       int64   `json:"stock" validate:"required,numeric"`
}

type UpdateProductRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`

	Id          string  `params:"id" validate:"required,uuid"`
	CategoryId  string  `json:"category_id" validate:"omitempty,uuid"`
	Name        string  `json:"name" validate:"required,max=255,min=3"`
	Description *string `json:"description" validate:"omitempty,max=255,min=3"`
	ImageUrl    *string `json:"image_url" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,numeric"`
	Stock       int64   `json:"stock" validate:"required,numeric"`
}

type UpdateProductStockRequest struct {
	Items []UpdateStock `validate:"required,dive,required"`
}

type UpdateStock struct {
	ProductId string `json:"product_id" validate:"required,uuid"`
	Stock     int64  `json:"stock" validate:"required,numeric"`
}

type UpsertProductResponse struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	ShopId      string    `json:"shop_id" db:"shop_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	ImageUrl    *string   `json:"image_url" db:"image_url"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type DeleteProductRequest struct {
	ProductId string `params:"product_id" validate:"required,uuid"`
	UserId    string `query:"user_id" validate:"required,uuid"`
}

type GetProductsRequest struct {
	ShopId        string `query:"shop_id" validate:"omitempty,uuid"`
	CategoryId    string `query:"category_id" validate:"omitempty,uuid"`
	Name          string `query:"name" validate:"omitempty,max=255,min=3"`
	PriceMinStr   string `query:"price_min" validate:"omitempty,numeric,gte=0"`
	PriceMaxStr   string `query:"price_max" validate:"omitempty,numeric,gte=0"`
	IsAvailable   bool   `query:"is_available"`
	ProductIdsStr string `query:"product_ids"`

	Page  int `query:"page" validate:"required,min=1"`
	Limit int `query:"limit" validate:"required,min=1,max=100"`

	PriceMin   float64
	PriceMax   float64
	ProductIds []string
}

func (r *GetProductsRequest) SetDefaults() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Limit < 1 {
		r.Limit = 10
	}

	if r.ProductIdsStr != "" {
		// split product ids string by comma
		ids := strings.Split(r.ProductIdsStr, ",")
		r.ProductIds = append(r.ProductIds, ids...)
	} else {
		r.ProductIds = make([]string, 0)
	}
}

func (r *GetProductsRequest) CostumValidation() (int, map[string][]string) {
	var (
		errors   = make(map[string][]string)
		err      error
		priceMin float64
		priceMax float64
	)

	if r.PriceMinStr != "" {
		priceMin, err = strconv.ParseFloat(r.PriceMinStr, 64)
		if err != nil {
			errors["price_min"] = append(errors["price_min"], "price_min must be a number.")
		}
		r.PriceMin = priceMin
	}

	if r.PriceMaxStr != "" {
		priceMax, err = strconv.ParseFloat(r.PriceMaxStr, 64)
		if err != nil {
			errors["price_max"] = append(errors["price_max"], "price_max must be a number.")
		}
		r.PriceMax = priceMax
	}

	if len(errors) > 0 {
		return 400, errors
	}

	errors = nil
	return 0, errors
}

type GetProductsResponse struct {
	Items []Product `json:"items"`
	Meta  Meta      `json:"meta"`
}

type Product struct {
	Id         string    `json:"id" db:"id"`
	CategoryId string    `json:"category_id" db:"category_id"`
	ShopId     string    `json:"shop_id" db:"shop_id"`
	Name       string    `json:"name" db:"name"`
	ImageUrl   *string   `json:"image_url" db:"image_url"`
	Price      float64   `json:"price" db:"price"`
	Stock      int       `json:"stock" db:"stock"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Meta struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

func (m *Meta) CountTotalPage() {
	if m.TotalData == 0 {
		m.TotalPage = 0
		return
	}

	m.TotalPage = m.TotalData / m.Limit
	if m.TotalData%m.Limit > 0 {
		m.TotalPage++
	}
}
