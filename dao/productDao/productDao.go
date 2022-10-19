package productDao

import (
	"errors"

	"github.com/paper-trade-chatbot/be-common/pagination"
	models "github.com/paper-trade-chatbot/be-product/models/databaseModels"
	"github.com/paper-trade-chatbot/be-proto/general"
	"google.golang.org/genproto/googleapis/cloud/common"

	"gorm.io/gorm"
)

const table = "product"

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	productType  []models.ProductType
	exchangeCode []string
	Status       int
	Display      int
	Offset       int
	Limit        int
}

// New a row
func New(tx *gorm.DB, model *models.ProductModel) (uint64, error) {

	err := tx.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return model.ID, nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*models.ProductModel, error) {

	result := &models.ProductModel{}
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Gets return records as raw-data-form
func Gets(tx *gorm.DB, query *QueryModel) ([]models.ProductModel, error) {
	result := make([]models.ProductModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.ProductModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetsWithPagination(tx *gorm.DB, query *QueryModel, paginate *general.Pagination) ([]models.ProductModel, *general.PaginationInfo, error) {

	var rows []models.ProductModel
	var count int = 0
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Count(&count).
		Scopes(paginateChain(paginate)).
		Scan(&rows).Error

	offset, _ := common.GetOffsetAndLimit(paginate)
	paginationInfo := pagination.SetProductPaginationDto(paginate.Page, paginate.PageSize, count, offset)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.ProductModel{}, paginationInfo, nil
	}

	if err != nil {
		return []models.ProductModel{}, nil, err
	}

	return rows, paginationInfo, nil
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(productTypeInScope(query.productType)).
			Scopes(exchangeCodeInScope(query.exchangeCode)).
			Scopes(statusEqualScope(query.Status)).
			Scopes(displayEqualScope(query.Display)).
			Scopes(offsetScope(query.Offset)).
			Scopes(limitScope(query.Limit))

	}
}

func paginateChain(paginate *general.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := common.GetOffsetAndLimit(paginate)
		return db.
			Scopes(offsetScope(offset)).
			Scopes(limitScope(limit))

	}
}

func productTypeInScope(productType []models.ProductType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(productType) > 0 {
			return db.Where(table+".type = ?", productType)
		}
		return db
	}
}

func exchangeCodeInScope(exchange []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(exchange) > 0 {
			return db.Where(table+".exchange_code = ?", exchange)
		}
		return db
	}
}

func statusEqualScope(status int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != 0 {
			return db.Where(table+".status = ?", status)
		}
		return db
	}
}

func displayEqualScope(display int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if display != 0 {
			return db.Where(table+".display = ?", display)
		}
		return db
	}
}

func limitScope(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

func offsetScope(offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Limit(offset)
		}
		return db
	}
}
