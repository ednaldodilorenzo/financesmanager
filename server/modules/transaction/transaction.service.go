package transaction

import (
	"errors"
	"io"
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
	PrepareFileImport(fileReader io.Reader, accountId uint32, date *time.Time, fileType string) ([]TransactionUploadSchema, error)
}

type TransactionServiceStruct struct {
	generic.GenericService[*model.Transaction]
	repository      TransactionRepository
	accountService  account.AccountService
	categoryService category.CategoryService
	parserFactory   *util.ParserFactory
}

func NewTransactionService(service generic.GenericService[*model.Transaction], repository TransactionRepository, accountService account.AccountService, categoryService category.CategoryService) TransactionService {
	parserFactory := util.NewParserFactory()

	return &TransactionServiceStruct{
		service,
		repository,
		accountService,
		categoryService,
		parserFactory,
	}
}

func (ts *TransactionServiceStruct) FindAllRelated(month *int, year *int) ([]model.Transaction, error) {
	return ts.repository.FindAllWithRelationships(month, year)
}

func (ts *TransactionServiceStruct) isDuplicated(value int32, paymentDate time.Time, transactionDate time.Time) (bool, error) {
	transaction, err := ts.repository.FindOneByValuePaymentDateAndTransactionDate(value, paymentDate, transactionDate)
	if err != nil {
		var runtimeError *util.RuntimeError
		if errors.As(err, &runtimeError) {
			return false, runtimeError
		}
	}

	duplicated := false
	if transaction != nil {
		duplicated = true
	}

	return duplicated, nil
}

func (ts *TransactionServiceStruct) PrepareFileImport(fileReader io.Reader, accountId uint32, date *time.Time, fileType string) ([]TransactionUploadSchema, error) {
	var constFileType util.FileImportType
	switch fileType {
	case "BBCA":
		constFileType = util.BBCA
	case "C6CC":
		constFileType = util.C6CC
	default:
		constFileType = util.CUAL
	}

	// Retrieve the parser from the factory
	parser, err := ts.parserFactory.GetParser(constFileType)
	if err != nil {
		return nil, err
	}
	parsedData, err := parser(fileReader, date)

	if err != nil {
		return nil, err
	}

	var transactionData []TransactionUploadSchema
	for _, record := range parsedData {
		duplicated, err := ts.isDuplicated(record.Value, record.PaymentDate, record.TransactionDate)
		if err != nil {
			return nil, err
		}

		var accountID uint32
		if record.AccountName != nil {
			account, err := ts.accountService.FindByName(*record.AccountName)
			if err != nil {
				return nil, err
			}

			accountID = account.ID
		} else {
			accountID = accountId
		}

		var categoryID uint32
		if record.CategoryName != nil {
			category, err := ts.categoryService.FindByName(*record.CategoryName)
			if err != nil {
				return nil, err
			}
			if category != nil {
				categoryID = category.ID
			}

		} else {
			categoryID = 0
		}

		newRecord := TransactionUploadSchema{
			CategoryId:      &categoryID,
			AccountId:       accountID,
			Description:     record.Description,
			Value:           record.Value,
			PaymentDate:     record.PaymentDate,
			TransactionDate: record.TransactionDate,
			Duplicated:      duplicated,
		}

		transactionData = append(transactionData, newRecord)
	}
	return transactionData, nil
}
