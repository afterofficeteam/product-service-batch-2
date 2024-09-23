package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error)
	GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error)
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error
	UpdateProductStock(ctx context.Context, req *entity.UpdateProductStockRequest) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error)
	GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error)
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error
	UpdateProductStock(ctx context.Context, req *entity.UpdateProductStockRequest) error

	IsShopOwner(ctx context.Context, userId, shopId string) (bool, error)
	IsProductOwner(ctx context.Context, userId, productId string) (bool, error)
}
