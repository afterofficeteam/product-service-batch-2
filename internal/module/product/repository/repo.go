package repository

import (
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"database/sql"

	"codebase-app/internal/module/product/entity"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ports.ProductRepository {
	return &productRepository{
		db,
	}
}

func (p *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error) {
	var (
		res entity.UpsertProductResponse
	)

	query := `
		INSERT INTO
			products (
				shop_id,
				category_id,
				name,
				description,
				image_url,
				price,
				stock
			)
			VALUES ( $1, $2, $3, $4, $5, $6, $7 )
			RETURNING
				id, shop_id, name, description, image_url, price, stock, created_at, updated_at
	`

	err := p.db.QueryRowxContext(ctx, query,
		req.ShopId,
		req.CategoryId,
		req.Name,
		req.Description,
		req.ImageUrl,
		req.Price,
		req.Stock,
	).StructScan(&res)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: CreateProduct failed")
		return res, err
	}

	res.UserId = req.UserId
	return res, nil
}

func (p *productRepository) GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Product
	}
	var (
		res  entity.GetProductsResponse
		data = make([]dao, 0)
		arg  = make(map[string]any)
	)
	res.Meta.Page = req.Page
	res.Meta.Limit = req.Limit

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			category_id,
			shop_id,
			name,
			image_url,
			price,
			stock,
			created_at,
			updated_at
		FROM
			products
		WHERE
			deleted_at IS NULL
	`

	if len(req.ProductIds) > 0 {
		var ids string
		for i, id := range req.ProductIds {
			ids += "'" + id + "'"
			if i < len(req.ProductIds)-1 {
				ids += ", "
			}
		}
		query += " AND id IN (" + ids + ")"
	}

	if req.ShopId != "" {
		query += " AND shop_id = :shop_id"
		arg["shop_id"] = req.ShopId
	}

	if req.CategoryId != "" {
		query += " AND category_id = :category_id"
		arg["category_id"] = req.CategoryId
	}

	if req.Name != "" {
		query += " AND name ILIKE '%' || :name || '%'"
		arg["name"] = req.Name
	}

	if req.PriceMinStr != "" {
		query += " AND price >= :price_min"
		arg["price_min"] = req.PriceMin
	}

	if req.PriceMaxStr != "" {
		query += " AND price <= :price_max"
		arg["price_max"] = req.PriceMax
	}

	if req.IsAvailable {
		query += " AND stock > 0"
	}

	query += `
		ORDER BY created_at DESC
		LIMIT :limit
		OFFSET :offset
	`
	arg["limit"] = req.Limit
	arg["offset"] = (req.Page - 1) * req.Limit

	nstmt, err := p.db.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: GetProducts failed")
		return res, err
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &data, arg)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: GetProducts failed")
		return res, err
	}

	for _, d := range data {
		res.Items = append(res.Items, entity.Product{
			Id:         d.Id,
			CategoryId: d.CategoryId,
			ShopId:     d.ShopId,
			Name:       d.Name,
			ImageUrl:   d.ImageUrl,
			Price:      d.Price,
			Stock:      d.Stock,
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
		})

		res.Meta.TotalData = d.TotalData
	}

	res.Meta.CountTotalPage()
	return res, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error) {
	var (
		res entity.UpsertProductResponse
	)

	query := `
		UPDATE
			products
		SET
			category_id = $1,
			name = $2,
			description = $3,
			image_url = $4,
			price = $5,
			stock = $6,
			updated_at = NOW()
		WHERE
			id = $7
			AND deleted_at IS NULL
		RETURNING
			id, shop_id, name, description, image_url, price, stock, created_at, updated_at
	`

	err := p.db.QueryRowxContext(ctx, query,
		req.CategoryId,
		req.Name,
		req.Description,
		req.ImageUrl,
		req.Price,
		req.Stock,
		req.Id,
	).StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Any("payload", req).Msg("repository: Product not found")
			return res, errmsg.NewCustomErrors(404, errmsg.WithMessage("Product not found"))
		}
		log.Error().Err(err).Any("payload", req).Msg("repository: UpdateProduct failed")
		return res, err
	}

	res.UserId = req.UserId
	return res, nil
}
func (p *productRepository) UpdateProductStock(ctx context.Context, req *entity.UpdateProductStockRequest) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: UpdateProductStock failed")
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE
			products
		SET
			stock = :stock
		WHERE
			id = :id
	`

	for _, item := range req.Items {
		arg := map[string]any{
			"id":    item.ProductId,
			"stock": item.Stock,
		}

		_, err = tx.NamedExec(query, arg)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository: UpdateProductStock failed")
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: UpdateProductStock failed")
		return err
	}

	return nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
	UPDATE products
		SET deleted_at = NOW()
	WHERE
		id = $1
	`

	_, err := p.db.ExecContext(ctx, query, req.ProductId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: DeleteProduct failed")
		return err
	}

	return nil
}

func (p *productRepository) IsShopOwner(ctx context.Context, userId, shopId string) (bool, error) {
	var (
		isOwner bool
		payload = struct {
			UserId string `json:"user_id"`
			ShopId string `json:"shop_id"`
		}{userId, shopId}
	)

	query := `
		SELECT
			EXISTS (
				SELECT 1
				FROM
					shops
				WHERE
					user_id = $1
					AND id = $2
					AND deleted_at IS NULL
			)
	`

	err := p.db.GetContext(ctx, &isOwner, query, userId, shopId)
	if err != nil {
		log.Error().Err(err).Any("payload", payload).Msg("repository: IsShopOwner failed")
		return isOwner, err
	}

	return isOwner, nil
}

func (p *productRepository) IsProductOwner(ctx context.Context, userId, productId string) (bool, error) {
	var (
		isOwner bool
		payload = struct {
			UserId    string `json:"user_id"`
			ProductId string `json:"product_id"`
		}{userId, productId}
	)

	query := `
		SELECT
			EXISTS (
				SELECT 1
				FROM
					products
				LEFT JOIN
					shops ON products.shop_id = shops.id
				WHERE
					shops.user_id = $1
					AND products.id = $2
			)
	`

	err := p.db.GetContext(ctx, &isOwner, query, userId, productId)
	if err != nil {
		log.Error().Err(err).Any("payload", payload).Msg("repository: IsProductOwner failed")
		return isOwner, err
	}

	return isOwner, nil
}
