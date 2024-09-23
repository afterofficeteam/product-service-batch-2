package service

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/rs/zerolog/log"
)

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(r ports.ProductRepository) ports.ProductService {
	return &productService{
		repo: r,
	}
}

func (p *productService) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error) {
	var res entity.UpsertProductResponse

	isShopOwner, err := p.repo.IsShopOwner(ctx, req.UserId, req.ShopId)
	if err != nil {
		return res, err
	}

	if !isShopOwner {
		log.Warn().Any("payload", req).Msg("service: User is not shop owner")
		return res, errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not shop owner"))
	}

	res, err = p.repo.CreateProduct(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (p *productService) GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error) {
	res, err := p.repo.GetProducts(ctx, req)
	if err != nil {
		return res, err
	}

	if len(res.Items) == 0 {
		log.Warn().Any("payload", req).Msg("service: Products not found")
		return res, errmsg.NewCustomErrors(404, errmsg.WithMessage("Products not found"))
	}

	return res, nil
}

func (p *productService) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error) {
	var res entity.UpsertProductResponse

	isProductOwner, err := p.repo.IsProductOwner(ctx, req.UserId, req.Id)
	if err != nil {
		return res, err
	}

	if !isProductOwner {
		log.Warn().Any("payload", req).Msg("service: User is not product owner")
		return res, errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))
	}

	res, err = p.repo.UpdateProduct(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (p *productService) UpdateProductStock(ctx context.Context, req *entity.UpdateProductStockRequest) error {
	return p.repo.UpdateProductStock(ctx, req)
}

func (p *productService) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	isProductOwner, err := p.repo.IsProductOwner(ctx, req.UserId, req.ProductId)
	if err != nil {
		return err
	}

	if !isProductOwner {
		log.Warn().Any("payload", req).Msg("service: User is not product owner")
		return errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))
	}

	return p.repo.DeleteProduct(ctx, req)
}
