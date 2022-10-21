package product

import (
	"context"
	"database/sql"

	"github.com/asaskevich/govalidator"
	common "github.com/paper-trade-chatbot/be-common"
	"github.com/paper-trade-chatbot/be-product/dao/exchangeDao"
	"github.com/paper-trade-chatbot/be-product/dao/productDao"
	"github.com/paper-trade-chatbot/be-product/database"
	"github.com/paper-trade-chatbot/be-product/logging"
	models "github.com/paper-trade-chatbot/be-product/models/databaseModels"
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
	db := database.GetDB()

	queryModel := &exchangeDao.QueryModel{}

	models, paginationInfo, err := exchangeDao.GetsWithPagination(db, queryModel, in.Pagination)
	if err != nil {
		return nil, err
	}

	exchanges := []*product.Exchange{}

	if len(models) == 0 {
		return &product.GetExchangesRes{
			Exchange:       exchanges,
			PaginationInfo: paginationInfo,
		}, nil
	}

	for _, m := range models {

		var openTime, closeTime *int64
		if m.OpenTime.Valid {
			openTimeObject := m.OpenTime.Time.Unix()
			openTime = &openTimeObject
		}
		if m.CloseTime.Valid {
			closeTimeObject := m.CloseTime.Time.Unix()
			closeTime = &closeTimeObject
		}
		e := &product.Exchange{
			Code:           m.Code,
			ProductType:    product.ProductType(m.ProductType),
			Name:           m.Name,
			Status:         product.Status(m.Status),
			Display:        product.Display(m.Display),
			CountryCode:    m.CountryCode,
			TimezoneOffset: float64(m.TimezoneOffset),
			OpenTime:       openTime,
			CloseTime:      closeTime,
			DaylightSaving: m.DaylightSaving,
			Location:       m.Location,
			ExchangeDay:    nil,
			ExceptionTime:  nil,
			CreatedAt:      m.CreatedAt.Unix(),
			UpdatedAt:      m.CreatedAt.Unix(),
		}
		exchanges = append(exchanges, e)
	}

	return &product.GetExchangesRes{
		Exchange:       exchanges,
		PaginationInfo: paginationInfo,
	}, nil
}

func (impl *ProductImpl) CreateProduct(ctx context.Context, in *product.CreateProductReq) (*product.CreateProductRes, error) {
	db := database.GetDB()

	logging.Info(ctx, "[CreateProduct] %s %s", in.ExchangeCode, in.Code)

	checkProductForm := struct {
		Type         int    `valid:"range(1|4)"`
		ExchangeCode string `valid:"required"`
		ProductCode  string `valid:"required"`
	}{
		Type:         int(in.Type),
		ExchangeCode: in.ExchangeCode,
		ProductCode:  in.Code,
	}

	if _, err := govalidator.ValidateStruct(checkProductForm); err != nil {
		logging.Info(ctx, "[CreateProduct] err: %v", err)
		return nil, common.ErrNoRequiredParam
	}

	if in.Status == 0 {
		in.Status = 1
	}

	if in.Display == 0 {
		in.Display = 1
	}

	var minimumOrder sql.NullFloat64
	if in.MinimumOrder != nil {
		minimumOrder.Valid = true
		minimumOrder.Float64 = in.GetMinimumOrder()
	}

	var iconID sql.NullString
	if in.IconID != nil {
		iconID.Valid = true
		iconID.String = in.GetIconID()
	}

	_, err := productDao.New(db, &models.ProductModel{
		Type:         models.ProductType(in.GetType()),
		ExchangeCode: in.GetExchangeCode(),
		Code:         in.GetCode(),
		Name:         in.GetName(),
		Status:       int(in.GetStatus()),
		Display:      int(in.GetDisplay()),
		CurrencyCode: in.GetCurrencyCode(),
		TickUnit:     in.GetTickUnit(),
		MinimumOrder: minimumOrder,
		IconID:       iconID,
	})
	if err != nil {
		return nil, err
	}

	return &product.CreateProductRes{}, nil
}

func (impl *ProductImpl) GetProduct(ctx context.Context, in *product.GetProductReq) (*product.GetProductRes, error) {
	db := database.GetDB()

	queryModel := &productDao.QueryModel{}

	switch query := in.GetProduct().(type) {
	case *product.GetProductReq_Id:
		queryModel.ID = uint64(query.Id)
	case *product.GetProductReq_Code:
		queryModel.ExchangeCode = query.Code.GetExchangeCode()
		queryModel.Code = query.Code.GetExchangeCode()
	}

	model, err := productDao.Get(db, queryModel)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return &product.GetProductRes{}, nil
	}

	var minimumOrder *float64
	if model.MinimumOrder.Valid {
		minimumOrderObject := model.MinimumOrder.Float64
		minimumOrder = &minimumOrderObject
	}
	var iconID *string
	if model.IconID.Valid {
		iconIDObject := model.IconID.String
		iconID = &iconIDObject
	}

	return &product.GetProductRes{
		Product: &product.Product{
			Id:           int64(model.ID),
			Type:         product.ProductType(model.Type),
			ExchangeCode: model.ExchangeCode,
			Code:         model.Code,
			Name:         model.Name,
			Status:       product.Status(model.Status),
			Display:      product.Display(model.Display),
			CurrencyCode: model.CurrencyCode,
			TickUnit:     model.TickUnit,
			MinimumOrder: minimumOrder,
			IconID:       iconID,
			CreatedAt:    model.CreatedAt.Unix(),
			UpdatedAt:    model.UpdatedAt.Unix(),
		},
	}, nil
}

func (impl *ProductImpl) GetProducts(ctx context.Context, in *product.GetProductsReq) (*product.GetProductsRes, error) {
	db := database.GetDB()

	queryModel := &productDao.QueryModel{
		ExchangeCodes: in.ExchangeCode,
	}

	for _, t := range in.ProductType {
		queryModel.ProductType = append(queryModel.ProductType, models.ProductType(t))
	}

	if in.Status != nil {
		queryModel.Status = int(in.GetStatus())
	}

	if in.Display != nil {
		queryModel.Display = int(in.GetDisplay())
	}

	models, paginationInfo, err := productDao.GetsWithPagination(db, queryModel, in.Pagination)
	if err != nil {
		return nil, err
	}

	products := []*product.Product{}

	if len(models) == 0 {
		return &product.GetProductsRes{
			Product:        products,
			PaginationInfo: paginationInfo,
		}, nil
	}

	for _, m := range models {

		var minimumOrder *float64
		if m.MinimumOrder.Valid {
			minimumOrderObject := m.MinimumOrder.Float64
			minimumOrder = &minimumOrderObject
		}
		var iconID *string
		if m.IconID.Valid {
			iconIDObject := m.IconID.String
			iconID = &iconIDObject
		}
		p := &product.Product{
			Id:           int64(m.ID),
			Type:         product.ProductType(m.Type),
			ExchangeCode: m.ExchangeCode,
			Code:         m.Code,
			Name:         m.Name,
			Status:       product.Status(m.Status),
			Display:      product.Display(m.Display),
			CurrencyCode: m.CurrencyCode,
			TickUnit:     m.TickUnit,
			MinimumOrder: minimumOrder,
			IconID:       iconID,
			CreatedAt:    m.CreatedAt.Unix(),
			UpdatedAt:    m.UpdatedAt.Unix(),
		}
		products = append(products, p)
	}

	return &product.GetProductsRes{
		Product:        products,
		PaginationInfo: paginationInfo,
	}, nil
}

func (impl *ProductImpl) ModifyProduct(ctx context.Context, in *product.ModifyProductReq) (*product.ModifyProductRes, error) {
	return nil, common.ErrNotImplemented
}

func (impl *ProductImpl) DeleteProduct(ctx context.Context, in *product.DeleteProductReq) (*product.DeleteProductRes, error) {
	return nil, common.ErrNotImplemented
}
