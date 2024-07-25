package transaction

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	generic.GenericRepository[*model.Transaction]
	FindAllWithRelationships(*int, *int) ([]model.Transaction, error)
}

type TransactionRespositoryStruct struct {
	generic.GenericRepository[*model.Transaction]
	db *gorm.DB
}

func NewTransactionRepository(repository generic.GenericRepository[*model.Transaction], database *config.Database) TransactionRepository {
	return &TransactionRespositoryStruct{
		repository,
		database.DB,
	}
}

func (tr *TransactionRespositoryStruct) FindAllWithRelationships(month *int, year *int) ([]model.Transaction, error) {
	var items []model.Transaction

	query := tr.db.Model(&items).Preload("Account").Preload("Category")

	if month != nil && year != nil {
		query = query.Where("MONTH(date) = ? AND YEAR(date) = ?", *month, *year)
	} else if month != nil {
		query = query.Where("MONTH(date) = ?", *month)
	} else if year != nil {
		query = query.Where("YEAR(date) = ?", *year)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
