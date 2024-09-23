package service

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	mockPort "codebase-app/mock/module/product/ports"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockService struct {
	mock.Mock
}

type ServiceList struct {
	suite.Suite
	mockProductRepo *mockPort.MockProductRepo
	service         ports.ProductService

	mockCreateProductReq          *entity.CreateProductRequest
	mockUpdateProductReq          *entity.UpdateProductRequest
	mockGetProductsReq            *entity.GetProductsRequest
	mockGetProductRes             entity.GetProductsResponse
	mockGetProductEmptyProductRes entity.GetProductsResponse
}

func (suite *ServiceList) SetupTest() {
	suite.mockProductRepo = new(mockPort.MockProductRepo)
	suite.service = NewProductService(suite.mockProductRepo)
	suite.mockCreateProductReq = &entity.CreateProductRequest{
		UserId:      "1",
		ShopId:      "2",
		CategoryId:  "3",
		Name:        "Product 1",
		Description: nil,
		ImageUrl:    nil,
		Price:       1000,
		Stock:       10,
	}

	suite.mockUpdateProductReq = &entity.UpdateProductRequest{
		UserId:      "1",
		Id:          "2",
		CategoryId:  "3",
		Name:        "Product 1",
		Description: nil,
		ImageUrl:    nil,
		Price:       1000,
		Stock:       10,
	}

	suite.mockGetProductsReq = &entity.GetProductsRequest{
		ShopId: "1",
	}

	suite.mockGetProductRes = entity.GetProductsResponse{
		Items: []entity.Product{
			{
				Id:         "1",
				CategoryId: "1",
				ShopId:     "1",
				Name:       "Product 1",
				ImageUrl:   nil,
				Price:      1000,
				Stock:      10,
			},
		},
		Meta: entity.Meta{
			TotalData: 1,
			TotalPage: 1,
			Page:      1,
			Limit:     10,
		},
	}

	suite.mockGetProductEmptyProductRes = entity.GetProductsResponse{
		Items: []entity.Product{},
		Meta: entity.Meta{
			TotalData: 0,
			TotalPage: 1,
			Page:      1,
			Limit:     10,
		},
	}
}

// Testing CreateProduct

func (u *ServiceList) TestCreateProduct_Success() {
	ctx := context.Background()
	req := u.mockCreateProductReq
	u.mockProductRepo.Mock.On("IsShopOwner", ctx, req.UserId, req.ShopId).Return(true, nil)
	u.mockProductRepo.Mock.On("CreateProduct", ctx, req).Return(mock.Anything, nil)
	_, err := u.service.CreateProduct(ctx, req)

	u.Equal(nil, err)
}

func (u *ServiceList) TestCreateProduct_IsShopOwnerError() {
	ctx := context.Background()
	req := u.mockCreateProductReq
	u.mockProductRepo.Mock.On("IsShopOwner", ctx, req.UserId, req.ShopId).Return(false, errors.New(mock.Anything))
	_, err := u.service.CreateProduct(ctx, req)

	u.Equal(errors.New(mock.Anything), err)
}

func (u *ServiceList) TestCreateProduct_UserIsNotTheShopOwner() {
	ctx := context.Background()
	req := u.mockCreateProductReq
	errForbidden := errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not shop owner"))

	u.mockProductRepo.Mock.On("IsShopOwner", ctx, req.UserId, req.ShopId).Return(false, nil)
	_, err := u.service.CreateProduct(ctx, req)

	u.Equal(errForbidden, err)
}
func (u *ServiceList) TestCreateProduct_Fail() {
	ctx := context.Background()
	req := u.mockCreateProductReq

	u.mockProductRepo.Mock.On("IsShopOwner", ctx, req.UserId, req.ShopId).Return(true, nil)
	u.mockProductRepo.Mock.On("CreateProduct", ctx, req).Return(mock.Anything, errors.New(mock.Anything))
	_, err := u.service.CreateProduct(ctx, req)

	u.Equal(errors.New(mock.Anything), err)
}

// Testing GetProducts

func (suite *ServiceList) TestGetProducts_Success() {
	ctx := context.Background()
	req := suite.mockGetProductsReq
	res := suite.mockGetProductRes

	suite.mockProductRepo.On("GetProducts", ctx, req).Return(res, nil)
	_, err := suite.service.GetProducts(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestGetProducts_GetProductsRepoError() {
	ctx := context.Background()
	reqMock := suite.mockGetProductsReq
	resMock := suite.mockGetProductRes

	suite.mockProductRepo.On("GetProducts", ctx, reqMock).Return(resMock, errors.New("error"))
	_, err := suite.service.GetProducts(ctx, reqMock)

	suite.Equal(errors.New("error"), err)
}

func (suite *ServiceList) TestGetProducts_ProductsEmpty() {
	ctx := context.Background()
	req := suite.mockGetProductsReq
	res := suite.mockGetProductEmptyProductRes
	errProductEmpty := errmsg.NewCustomErrors(404, errmsg.WithMessage("Products not found"))

	suite.mockProductRepo.On("GetProducts", ctx, req).Return(res, nil)
	_, err := suite.service.GetProducts(ctx, req)

	suite.Equal(errProductEmpty, err)
}

// Testing UpdateProduct

func (suite *ServiceList) TestUpdateProduct_Success() {
	ctx := context.Background()
	reqMock := &entity.UpdateProductRequest{
		UserId: "1",
		Id:     "1",
	}

	resMock := entity.UpsertProductResponse{
		Id: "1",
	}

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.Id).Return(true, nil)
	suite.mockProductRepo.On("UpdateProduct", ctx, reqMock).Return(resMock, nil)
	_, err := suite.service.UpdateProduct(ctx, reqMock)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestUpdateProduct_IsProductOwnerError() {
	ctx := context.Background()
	reqMock := &entity.UpdateProductRequest{
		UserId: "1",
		Id:     "1",
	}

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.Id).Return(false, errors.New("error"))
	_, err := suite.service.UpdateProduct(ctx, reqMock)

	suite.Equal(errors.New("error"), err)
}

func (suite *ServiceList) TestUpdateProduct_UserIsNotTheProductOwner() {
	ctx := context.Background()
	reqMock := &entity.UpdateProductRequest{
		UserId: "1",
		Id:     "1",
	}

	errForbidden := errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.Id).Return(false, nil)
	_, err := suite.service.UpdateProduct(ctx, reqMock)

	suite.Equal(errForbidden, err)
}

func (suite *ServiceList) TestUpdateProduct_UpdateProductError() {
	ctx := context.Background()
	reqMock := &entity.UpdateProductRequest{
		UserId: "1",
		Id:     "1",
	}

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.Id).Return(true, nil)
	suite.mockProductRepo.On("UpdateProduct", ctx, reqMock).Return(mock.Anything, errors.New(mock.Anything))
	_, err := suite.service.UpdateProduct(ctx, reqMock)

	suite.Equal(errors.New(mock.Anything), err)
}

// Testing UpdateProductStock

func (suite *ServiceList) TestUpdateProductStock_Success() {
	ctx := context.Background()
	reqMock := &entity.UpdateProductStockRequest{
		Items: []entity.UpdateStock{},
	}

	suite.mockProductRepo.On("UpdateProductStock", ctx, reqMock).Return(nil)
	err := suite.service.UpdateProductStock(ctx, reqMock)

	suite.Equal(nil, err)
}

// Testing DeleteProduct
func (suite *ServiceList) TestDeleteProduct_Success() {
	ctx := context.Background()
	reqMock := &entity.DeleteProductRequest{
		UserId:    "1",
		ProductId: "1",
	}

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.ProductId).Return(true, nil)
	suite.mockProductRepo.On("DeleteProduct", ctx, reqMock).Return(nil)
	err := suite.service.DeleteProduct(ctx, reqMock)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestDeleteProduct_IsProductOwnerError() {
	ctx := context.Background()
	reqMock := &entity.DeleteProductRequest{
		UserId:    "1",
		ProductId: "1",
	}

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.ProductId).Return(false, errors.New(mock.Anything))
	err := suite.service.DeleteProduct(ctx, reqMock)

	suite.Equal(errors.New(mock.Anything), err)
}

func (suite *ServiceList) TestDeleteProduct_UserIsNotTheProductOwner() {
	ctx := context.Background()
	reqMock := &entity.DeleteProductRequest{
		UserId:    "1",
		ProductId: "1",
	}

	errForbidden := errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))

	suite.mockProductRepo.On("IsProductOwner", ctx, reqMock.UserId, reqMock.ProductId).Return(false, nil)
	err := suite.service.DeleteProduct(ctx, reqMock)

	suite.Equal(errForbidden, err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceList))
}
