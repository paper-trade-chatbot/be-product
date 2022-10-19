package exchangeDao

import (
	"errors"

	models "github.com/paper-trade-chatbot/be-product/models/databaseModels"

	"gorm.io/gorm"
)

const table = "exchange"

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	Code    string
	Status  int
	Display int
}

// New a row
func New(tx *gorm.DB, model *models.ExchangeModel) (string, error) {

	err := tx.Table(table).
		Create(model).Error

	if err != nil {
		return "", err
	}
	return model.Code, nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*models.ExchangeModel, error) {

	result := &models.ExchangeModel{}
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
func Gets(tx *gorm.DB, query *QueryModel) ([]models.ExchangeModel, error) {
	result := make([]models.ExchangeModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.ExchangeModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(codeEqualScope(query.Code)).
			Scopes(statusEqualScope(query.Status)).
			Scopes(displayEqualScope(query.Display))

	}
}

func codeEqualScope(code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if code != "" {
			return db.Where(table+".code = ?", code)
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
