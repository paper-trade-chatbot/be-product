package permissionpathdao

import (
	"errors"

	models "github.com/paper-trade-chatbot/be-product/models/databaseModels"

	"gorm.io/gorm"
)

const table = "group_permission_path"

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	GroupID uint64
	Method  string
	Path    string
}

// New a row
func New(tx *gorm.DB, model *models.GroupPermissionPathModel) (uint64, error) {

	err := tx.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return model.ID, nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*models.GroupPermissionPathModel, error) {

	result := &models.GroupPermissionPathModel{}
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
func Gets(tx *gorm.DB, query *QueryModel) ([]models.GroupPermissionPathModel, error) {
	result := make([]models.GroupPermissionPathModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.GroupPermissionPathModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(groupIDEqualScope(query.GroupID)).
			Scopes(methodEqualScope(query.Method)).
			Scopes(pathEqualScope(query.Path))

	}
}

func groupIDEqualScope(groupID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if groupID != 0 {
			return db.Where(table+".group_id = ?", groupID)
		}
		return db
	}
}

func methodEqualScope(method string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if method != "" {
			return db.Where(table+".method = ?", method)
		}
		return db
	}
}

func pathEqualScope(path string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if path != "" {
			return db.Where(table+".path = ?", path)
		}
		return db
	}

}
