package category

import (
	"errors"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	generic.GenericRepository[*model.Category]
	FindByName(name string) (*model.Category, error)
}

type CategoryRespositoryStruct struct {
	generic.GenericRepository[*model.Category]
	dbConfig *config.Database
}

func NewAccountRepository(repository generic.GenericRepository[*model.Category], database *config.Database) CategoryRepository {
	return &CategoryRespositoryStruct{
		repository,
		database,
	}
}

func (cr *CategoryRespositoryStruct) FindByName(name string) (*model.Category, error) {
	var result model.Category

	err := cr.dbConfig.DB.First(&result, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle the case where no record is found
			return nil, nil
		}
		// Handle other errors (e.g., database connection issues)
		return nil, util.NewRuntimeError("Error in database operation", err)
	}

	return &result, nil
}
