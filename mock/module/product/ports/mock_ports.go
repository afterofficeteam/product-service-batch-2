package mock_ports

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockProductRepo struct {
	mock.Mock
}

func NewMockProductRepo() *MockProductRepo {
	return &MockProductRepo{}
}

var _ ports.ProductRepository = &MockProductRepo{}

func (m *MockProductRepo) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error) {
	args := m.Called(ctx, req)
	var (
		resp entity.UpsertProductResponse
		err  error
	)

	if n, ok := args.Get(0).(entity.UpsertProductResponse); ok {

		resp = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return resp, err
}

func (m *MockProductRepo) GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error) {
	args := m.Called(ctx, req)
	var (
		resp entity.GetProductsResponse
		err  error
	)

	if n, ok := args.Get(0).(entity.GetProductsResponse); ok {

		resp = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return resp, err
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error) {
	args := m.Called(ctx, req)
	var (
		resp entity.UpsertProductResponse
		err  error
	)

	if n, ok := args.Get(0).(entity.UpsertProductResponse); ok {

		resp = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return resp, err
}

func (m *MockProductRepo) UpdateProductStock(ctx context.Context, req *entity.UpdateProductStockRequest) error {
	args := m.Called(ctx, req)
	var (
		err error
	)

	if n, ok := args.Get(0).(error); ok {

		err = n
	}

	return err
}

func (m *MockProductRepo) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	args := m.Called(ctx, req)
	var (
		err error
	)

	if n, ok := args.Get(0).(error); ok {

		err = n
	}

	return err
}

func (m *MockProductRepo) IsShopOwner(ctx context.Context, userId, shopId string) (bool, error) {
	args := m.Called(ctx, userId, shopId)
	var (
		resp bool
		err  error
	)

	if n, ok := args.Get(0).(bool); ok {

		resp = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return resp, err
}

func (m *MockProductRepo) IsProductOwner(ctx context.Context, userId, productId string) (bool, error) {
	args := m.Called(ctx, userId, productId)
	var (
		resp bool
		err  error
	)

	if n, ok := args.Get(0).(bool); ok {

		resp = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return resp, err
}
