package product

import (
	"context"

	common "github.com/paper-trade-chatbot/be-common"
	"github.com/paper-trade-chatbot/be-product/dao/exchangeDao"
	"github.com/paper-trade-chatbot/be-product/database"
	"github.com/paper-trade-chatbot/be-proto/product"
)

type ProductIntf interface {
	GetExchange(ctx context.Context, in *product.GetExchangeReq) (*product.GetExchangeRes, error)
	GetExchanges(ctx context.Context, in *product.GetExchangesReq) (*product.GetExchangesRes, error)
	CreateProduct(ctx context.Context, in *product.CreateProductReq) (*product.CreateProductRes, error)
	GetProduct(ctx context.Context, in *product.GetProductReq) (*product.GetProductRes, error)
	GetProducts(ctx context.Context, in *product.GetProductsReq) (*product.GetProductsRes, error)
	ModifyProduct(ctx context.Context, in *product.ModifyProductReq) (*product.ModifyProductRes, error)
	DeleteProduct(ctx context.Context, in *product.DeleteProductReq) (*product.DeleteProductRes, error)
}

type ProductImpl struct {
	ProductClient product.ProductServiceClient
}

func New() ProductIntf {
	return &ProductImpl{}
}

func (impl *ProductImpl) GetExchange(ctx context.Context, in *product.GetExchangeReq) (*product.GetExchangeRes, error) {
	db := database.GetDB()

	queryModel := &exchangeDao.QueryModel{
		Code:    in.GetCode(),
		Status:  int(in.GetStatus()),
		Display: int(in.GetDisplay()),
	}

	model, err := exchangeDao.Get(db, queryModel)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return &product.GetExchangeRes{}, nil
	}

	var openTime, closeTime *int64
	if model.OpenTime.Valid {
		openTimeObject := model.OpenTime.Time.Unix()
		openTime = &openTimeObject
	}
	if model.CloseTime.Valid {
		closeTimeObject := model.CloseTime.Time.Unix()
		closeTime = &closeTimeObject
	}

	return &product.GetExchangeRes{
		Exchange: &product.Exchange{
			Code:           model.Code,
			ProductType:    product.ProductType(model.ProductType),
			Name:           model.Name,
			Status:         product.Status(model.Status),
			Display:        product.Display(model.Display),
			CountryCode:    model.CountryCode,
			TimezoneOffset: float64(model.TimezoneOffset),
			OpenTime:       openTime,
			CloseTime:      closeTime,
			DaylightSaving: model.DaylightSaving,
			Location:       model.Location,
			ExchangeDay:    nil,
			ExceptionTime:  nil,
			CreatedAt:      model.CreatedAt.Unix(),
			UpdatedAt:      model.UpdatedAt.Unix(),
		},
	}, nil

}

func (impl *ProductImpl) GetExchanges(ctx context.Context, in *product.GetExchangesReq) (*product.GetExchangesRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) CreateProduct(ctx context.Context, in *product.CreateProductReq) (*product.CreateProductRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) GetProduct(ctx context.Context, in *product.GetProductReq) (*product.GetProductRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) GetProducts(ctx context.Context, in *product.GetProductsReq) (*product.GetProductsRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) ModifyProduct(ctx context.Context, in *product.ModifyProductReq) (*product.ModifyProductRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) DeleteProduct(ctx context.Context, in *product.DeleteProductReq) (*product.DeleteProductRes, error) {
	return nil, common.ErrNotImplemented
}
