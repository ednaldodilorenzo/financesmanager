package account

import (
	"errors"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"gorm.io/gorm"
)

type AccountRepository interface {
	generic.GenericRepository[*model.Account]
	FindByName(name string) (*model.Account, error)
}

type AccountRespositoryStruct struct {
	generic.GenericRepository[*model.Account]
	dbConfig *config.Database
}

func NewAccountRepository(repository generic.GenericRepository[*model.Account], database *config.Database) AccountRepository {
	return &AccountRespositoryStruct{
		repository,
		database,
	}
}

func (ar *AccountRespositoryStruct) FindByName(name string) (*model.Account, error) {
	var result model.Account

	err := ar.dbConfig.DB.First(&result, "name = ?", name).Error
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
