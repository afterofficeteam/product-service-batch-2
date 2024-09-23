package rest

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"encoding/json"
	"time"

	llog "log"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type producthandler struct {
	service ports.ProductService
}

func NewProductHandler() *producthandler {
	repo := repository.NewProductRepository(adapter.Adapters.ShopeefunPostgres)
	service := service.NewProductService(repo)

	return &producthandler{
		service: service,
	}
}

func (h *producthandler) Register(router fiber.Router) {
	router.Get("/products", h.getProducts)

	router.Post("/products", m.UserIdHeader, h.createProduct)
	router.Patch("/product-stocks", h.updateProductStock)
	router.Patch("/products/:id", m.UserIdHeader, h.updateProduct)
	router.Delete("/products/:id", m.UserIdHeader, h.deleteProduct)
}

func (h *producthandler) createProduct(c *fiber.Ctx) error {
	var (
		req = &entity.CreateProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *producthandler) updateProduct(c *fiber.Ctx) error {
	var (
		req = &entity.UpdateProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = c.Query("user_id")
	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *producthandler) updateProductStock(c *fiber.Ctx) error {
	var (
		reqArray  = make([]entity.UpdateStock, 0)
		ctx       = c.Context()
		v         = adapter.Adapters.Validator
		req       = &entity.UpdateProductStockRequest{}
		message   = "Your request has been successfully processed"
		timeStart = time.Now()
	)
	defer func() {
		duration := time.Since(timeStart).Milliseconds()
		llog.Println("create payment request", duration)
		log.Debug().Any("duration", duration).Any("unit", "ms").Msg("service: UpdateProductStock")
	}()

	err := json.Unmarshal(c.Body(), &reqArray)
	if err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request body")
		// return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	req.Items = reqArray
	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		// code, errs := errmsg.Errors(err, req)
		// return c.Status(code).JSON(response.Error(errs))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = h.service.UpdateProductStock(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err, nil)
		response := response.Error(errs)
		message, ok := response["message"].(string)
		if !ok {
			message = "Your request has been failed to process"
		}
		return c.Status(code).SendString(message)
	}

	return c.Status(fiber.StatusOK).SendString(message)
}

func (h *producthandler) deleteProduct(c *fiber.Ctx) error {
	var (
		req = &entity.DeleteProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Query("user_id")
	req.ProductId = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *producthandler) getProducts(c *fiber.Ctx) error {
	var (
		req = &entity.GetProductsRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefaults()

	if code, errs := req.CostumValidation(); code != 0 {
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProducts(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
