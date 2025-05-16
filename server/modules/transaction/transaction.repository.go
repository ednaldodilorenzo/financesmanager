package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	generic.GenericRepository[*model.Transaction]
	FindAllWithRelationships(*int, *int, int) ([]model.Transaction, error)
	FindOneByValuePaymentDateAndTransactionDate(value int32, paymentDate time.Time, transactionDate time.Time, userId int) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, db *gorm.DB, item model.Transaction) error
}

type transactionRespository struct {
	generic.GenericRepository[*model.Transaction]
	dbConfig *config.Database
}

func NewTransactionRepository(repository generic.GenericRepository[*model.Transaction], database *config.Database) TransactionRepository {
	return &transactionRespository{
		GenericRepository: repository,
		dbConfig:          database,
	}
}

func (g *transactionRespository) CreateTransaction(ctx context.Context, db *gorm.DB, item model.Transaction) error {

	if err := db.Clauses(clause.Returning{}).Create(&item).Error; err != nil {
		return err
	}

	return nil
}

func (g *transactionRespository) FindById(id, userId int) (**model.Transaction, error) {
	var item *model.Transaction

	// Use Joins to enforce INNER JOIN
	err := g.dbConfig.DB.Model(&item).
		Joins("INNER JOIN category ON transaction.category_id = category.id AND transaction.user_id = category.user_id").
		Joins("INNER JOIN account ON transaction.account_id = account.id AND transaction.user_id = account.user_id").
		Joins("LEFT JOIN transaction_tag ON transaction.id = transaction_tag.transaction_id AND transaction.user_id = transaction_tag.user_id").
		Preload("Tags"). // Still preload tags to attach them to the struct
		Preload("Category").
		Preload("Account").
		First(&item, "transaction.id = ? AND transaction.user_id = ?", id, userId).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (tr *transactionRespository) FindAllWithRelationships(month *int, year *int, userId int) ([]model.Transaction, error) {
	var items []model.Transaction

	query := tr.dbConfig.DB.Model(&items).Preload("Account").Preload("Category")

	if month != nil && year != nil {
		query = query.Where("EXTRACT(MONTH FROM payment_date) = ? AND EXTRACT(YEAR FROM payment_date) = ? AND user_id = ? ORDER BY payment_date DESC, transaction_date DESC", *month, *year, userId)
	} else if month != nil {
		query = query.Where("EXTRACT(MONTH FROM payment_date) = ? AND user_id = ? ORDER BY payment_date DESC, transaction_date DESC", *month, userId)
	} else if year != nil {
		query = query.Where("EXTRACT(YEAR FROM payment_date) = ? AND user_id = ? ORDER BY payment_date DESC, transaction_date DESC", *year, userId)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (tr *transactionRespository) FindOneByValuePaymentDateAndTransactionDate(value int32, paymentDate time.Time, transactionDate time.Time, userId int) (*model.Transaction, error) {
	var result model.Transaction

	err := tr.dbConfig.DB.First(&result, "value = ? AND payment_date = ? AND transaction_date = ? AND user_id = ?", value, paymentDate, transactionDate, userId).Error
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
