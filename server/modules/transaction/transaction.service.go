package transaction

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/account"
	"github.com/ednaldo-dilorenzo/iappointment/modules/category"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type TransactionService interface {
	generic.GenericService[*model.Transaction]
	FindAllRelated(*int, *int) ([]model.Transaction, error)
	PrepareFileImport(fileReader io.Reader) ([]TransactionUploadSchema, error)
}

type TransactionServiceStruct struct {
	generic.GenericService[*model.Transaction]
	repository      TransactionRepository
	accountService  account.AccountService
	categoryService category.CategoryService
}

func NewTransactionService(service generic.GenericService[*model.Transaction], repository TransactionRepository, accountService account.AccountService, categoryService category.CategoryService) TransactionService {
	return &TransactionServiceStruct{
		service,
		repository,
		accountService,
		categoryService,
	}
}

func (ts *TransactionServiceStruct) FindAllRelated(month *int, year *int) ([]model.Transaction, error) {
	return ts.repository.FindAllWithRelationships(month, year)
}

func (ts *TransactionServiceStruct) PrepareFileImport(fileReader io.Reader) ([]TransactionUploadSchema, error) {
	// Parse the CSV file
	csvReader := csv.NewReader(fileReader)
	csvReader.Comma = '\t'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []TransactionUploadSchema
	// Iterate through the records and save them to the database
	for _, record := range records {

		value, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, err
		}

		paymentDate, err := time.Parse("02/01/2006", record[0]) // Adjust format as needed
		if err != nil {
			return nil, err
		}

		transactionDate, err := time.Parse("02/01/2006", record[5]) // Adjust format as needed
		if err != nil {
			return nil, err
		}

		transaction, err := ts.repository.FindOneByValueAndPaymentDate(float32(value*100), paymentDate)
		if err != nil {
			var runtimeError *util.RuntimeError
			if errors.As(err, &runtimeError) {
				return nil, runtimeError
			}
		}

		duplicated := false
		if transaction != nil {
			duplicated = true
		}

		account, err := ts.accountService.FindByName(record[4])
		if err != nil {
			var runtimeError *util.RuntimeError
			if errors.As(err, &runtimeError) {
				return nil, runtimeError
			}
		}

		category, err := ts.categoryService.FindByName(record[3])
		if err != nil {
			var runtimeError *util.RuntimeError
			if errors.As(err, &runtimeError) {
				return nil, runtimeError
			}
		}

		var categoryId uint32

		if category == nil {
			categoryId = 0
		} else {
			categoryId = category.ID
		}

		// Create a new Record instance
		dbRecord := TransactionUploadSchema{
			Description:     record[1],
			AccountId:       account.ID,
			CategoryId:      categoryId,
			Value:           int32(value * 100),
			PaymentDate:     paymentDate,
			TransactionDate: transactionDate,
			Duplicated:      duplicated,
		}

		transactions = append(transactions, dbRecord)

	}

	return transactions, nil
}
