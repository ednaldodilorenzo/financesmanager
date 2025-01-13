package transaction

import (
	"errors"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	generic.GenericRepository[*model.Transaction]
	FindById(id int) (**model.Transaction, error)
	FindAllWithRelationships(*int, *int) ([]model.Transaction, error)
	FindOneByValuePaymentDateAndTransactionDate(value int32, paymentDate time.Time, transactionDate time.Time) (*model.Transaction, error)
}

type TransactionRespositoryStruct struct {
	generic.GenericRepository[*model.Transaction]
	dbConfig *config.Database
}

func NewTransactionRepository(repository generic.GenericRepository[*model.Transaction], database *config.Database) TransactionRepository {
	return &TransactionRespositoryStruct{
		GenericRepository: repository,
		dbConfig:          database,
	}
}

func (g *TransactionRespositoryStruct) FindById(id int) (**model.Transaction, error) {
	var item *model.Transaction

	// Use Joins to enforce INNER JOIN
	err := g.dbConfig.DB.Model(&item).
		Joins("INNER JOIN category ON transaction.id_category = category.id AND transaction.id_user = category.id_user").
		Joins("INNER JOIN account ON transaction.id_account = account.id AND transaction.id_user = account.id_user").
		Joins("LEFT JOIN transaction_tag ON transaction.id = transaction_tag.transaction_id AND transaction.id_user = transaction_tag.user_id").
		Preload("Tags"). // Still preload tags to attach them to the struct
		Preload("Category").
		Preload("Account").
		First(&item, "transaction.id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (tr *TransactionRespositoryStruct) FindAllWithRelationships(month *int, year *int) ([]model.Transaction, error) {
	var items []model.Transaction

	query := tr.dbConfig.DB.Model(&items).Preload("Account").Preload("Category")

	if month != nil && year != nil {
		query = query.Where("EXTRACT(MONTH FROM payment_date) = ? AND EXTRACT(YEAR FROM payment_date) = ? ORDER BY payment_date DESC", *month, *year)
	} else if month != nil {
		query = query.Where("EXTRACT(MONTH FROM payment_date) = ? ORDER BY payment_date DESC", *month)
	} else if year != nil {
		query = query.Where("EXTRACT(YEAR FROM payment_date) = ? ORDER BY payment_date DESC", *year)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (tr *TransactionRespositoryStruct) FindOneByValuePaymentDateAndTransactionDate(value int32, paymentDate time.Time, transactionDate time.Time) (*model.Transaction, error) {
	var result model.Transaction

	err := tr.dbConfig.DB.First(&result, "value = ? AND payment_date = ? AND transaction_date = ?", value, paymentDate, transactionDate).Error
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
