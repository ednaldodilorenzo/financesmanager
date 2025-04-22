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
	FindAllRelated(*int, *int, int) ([]model.Transaction, error)
	PrepareFileImport(fileReader io.Reader, accountId uint32, date *time.Time, fileType string, userId int) ([]TransactionUploadSchema, error)
	CreateTransaction(*TransactionPostRequest, int) error
	UpdateTransaction(id int, item *TransactionPostRequest, userId int) error
}

type transactionService struct {
	generic.GenericService[*model.Transaction]
	repository      TransactionRepository
	accountService  account.AccountService
	categoryService category.CategoryService
	parserFactory   *util.ParserFactory
}

func NewTransactionService(service generic.GenericService[*model.Transaction], repository TransactionRepository, accountService account.AccountService, categoryService category.CategoryService) TransactionService {
	parserFactory := util.NewParserFactory()

	return &transactionService{
		service,
		repository,
		accountService,
		categoryService,
		parserFactory,
	}
}

func (ts *transactionService) CreateTransaction(transactionRequest *TransactionPostRequest, userId int) error {
	account, err := ts.accountService.FindById(transactionRequest.AccountId, userId)

	if err != nil {
		return err
	}

	if account == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"Account not found"})
	}

	var definedPaymentDate time.Time

	if (*account).Type == "C" {
		definedPaymentDate = time.Date(transactionRequest.PaymentYear, time.Month(transactionRequest.PaymentMonth), (*account).DueDay, 0, 0, 0, 0, time.UTC)
	} else {
		definedPaymentDate = transactionRequest.TransactionDate
	}

	newTransaction := &model.Transaction{
		CategoryId:      transactionRequest.CategoryId,
		AccountId:       transactionRequest.AccountId,
		Description:     transactionRequest.Description,
		Value:           transactionRequest.Value,
		PaymentDate:     definedPaymentDate,
		TransactionDate: transactionRequest.TransactionDate,
		Detail:          transactionRequest.Detail,
		Tags:            transactionRequest.Tags,
		UserDependent:   model.UserDependent{UserId: uint64(userId)},
	}

	err = ts.repository.Create(&newTransaction)

	if err != nil {
		return err
	}

	return nil
}

func (ts *transactionService) UpdateTransaction(id int, item *TransactionPostRequest, userId int) error {
	_, err := ts.repository.FindById(id, userId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.NewAPIError(util.ErrNotFound, []string{"Transaction not found"})
		}
	}

	account, err := ts.accountService.FindById(item.AccountId, userId)

	if err != nil {
		return err
	}

	if account == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"Account not found"})
	}

	var definedPaymentDate time.Time

	if (*account).Type == "C" {
		definedPaymentDate = time.Date(item.PaymentYear, time.Month(item.PaymentMonth), (*account).DueDay, 0, 0, 0, 0, time.UTC)
	} else {
		definedPaymentDate = (*item).TransactionDate
	}

	updatedTransaction := &model.Transaction{
		CategoryId:      item.CategoryId,
		AccountId:       item.AccountId,
		Description:     item.Description,
		Value:           item.Value,
		PaymentDate:     definedPaymentDate,
		TransactionDate: item.TransactionDate,
		Detail:          item.Detail,
		Tags:            item.Tags,
		UserDependent:   model.UserDependent{UserId: uint64(userId)},
	}

	if err := ts.repository.Update(id, &updatedTransaction, userId); err != nil {
		return err
	}

	return nil
}

func (ts *transactionService) FindAllRelated(month *int, year *int, userId int) ([]model.Transaction, error) {
	return ts.repository.FindAllWithRelationships(month, year, userId)
}

func (ts *transactionService) FindById(id int, userId int) (**model.Transaction, error) {
	return ts.repository.FindById(id, userId)
}

func (ts *transactionService) isDuplicated(value int32, paymentDate time.Time, transactionDate time.Time, userId int) (bool, error) {
	transaction, err := ts.repository.FindOneByValuePaymentDateAndTransactionDate(value, paymentDate, transactionDate, userId)
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

func (ts *transactionService) PrepareFileImport(fileReader io.Reader, accountId uint32, date *time.Time, fileType string, userId int) ([]TransactionUploadSchema, error) {
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
		duplicated, err := ts.isDuplicated(record.Value, record.PaymentDate, record.TransactionDate, userId)
		if err != nil {
			return nil, err
		}

		var accountID uint32
		if record.AccountName != nil {
			account, err := ts.accountService.FindByName(*record.AccountName, userId)
			if err != nil {
				return nil, err
			}

			accountID = account.Id
		} else {
			accountID = accountId
		}

		var categoryID uint32
		if record.CategoryName != nil {
			category, err := ts.categoryService.FindByName(*record.CategoryName, userId)
			if err != nil {
				return nil, err
			}
			if category != nil {
				categoryID = category.Id
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
