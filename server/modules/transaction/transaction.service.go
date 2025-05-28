package transaction

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/account"
	"github.com/ednaldo-dilorenzo/iappointment/modules/category"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type TransactionService interface {
	generic.GenericService[*model.Transaction]
	FindAllRelated(context.Context, *int, *int, int) ([]model.Transaction, error)
	PrepareFileImport(ctx context.Context, fileReader io.Reader, accountId uint32, paymentMonth uint8, paymentYear uint16, fileType string, userId int) ([]TransactionUploadSchema, error)
	CreateTransaction(ctx context.Context, transactionRequest *TransactionPostRequest, userId int) error
	UpdateTransaction(ctx context.Context, id int, item *TransactionPostRequest, userId int) error
}

type transactionService struct {
	generic.GenericService[*model.Transaction]
	repository      TransactionRepository
	accountService  account.AccountService
	categoryService category.CategoryService
	parserFactory   *util.ParserFactory
	txManager       config.TxManager
}

func NewTransactionService(service generic.GenericService[*model.Transaction], repository TransactionRepository, accountService account.AccountService, categoryService category.CategoryService, txManager config.TxManager) TransactionService {
	parserFactory := util.NewParserFactory()

	return &transactionService{
		service,
		repository,
		accountService,
		categoryService,
		parserFactory,
		txManager,
	}
}

func (ts *transactionService) CreateTransaction(ctx context.Context, transactionRequest *TransactionPostRequest, userId int) error {
	tx, err := ts.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	account, err := ts.accountService.FindById(ctx, transactionRequest.AccountId, userId)

	if err != nil {
		return err
	}

	if account == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"Account not found"})
	}

	var definedPaymentDate time.Time

	if (*account).Type == "C" {
		if transactionRequest.PaymentYear == 0 || transactionRequest.PaymentMonth == 0 {
			return util.NewAPIError(util.ErrBadRequest, []string{"Payment year and month are required"})
		}
		definedPaymentDate = time.Date(transactionRequest.PaymentYear, time.Month(transactionRequest.PaymentMonth), (*account).DueDay, 0, 0, 0, 0, time.UTC)
	} else {
		definedPaymentDate = transactionRequest.TransactionDate
	}

	newTransaction := model.Transaction{
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

	err = ts.repository.CreateTransaction(ctx, gormTx, newTransaction)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ts *transactionService) UpdateTransaction(ctx context.Context, id int, item *TransactionPostRequest, userId int) error {
	tx, err := ts.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	_, err = ts.repository.FindById(ctx, id, userId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.NewAPIError(util.ErrNotFound, []string{"Transaction not found"})
		}
	}

	account, err := ts.accountService.FindById(ctx, item.AccountId, userId)

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
		UserDependent:   model.UserDependent{Id: uint32(id), UserId: uint64(userId)},
	}

	err = ts.repository.Update(ctx, gormTx, id, &updatedTransaction, userId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ts *transactionService) FindAllRelated(ctx context.Context, month *int, year *int, userId int) ([]model.Transaction, error) {
	return ts.repository.FindAllWithRelationships(month, year, userId)
}

func (ts *transactionService) FindById(ctx context.Context, id int, userId int) (**model.Transaction, error) {
	return ts.repository.FindById(ctx, id, userId)
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

func (ts *transactionService) PrepareFileImport(ctx context.Context, fileReader io.Reader, accountId uint32, paymentMonth uint8, paymentYear uint16, fileType string, userId int) ([]TransactionUploadSchema, error) {
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

	account, err := ts.accountService.FindById(ctx, int(accountId), userId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, util.NewAPIError(util.ErrNotFound, []string{"Account not found"})
	}

	var date time.Time

	if (*account).Type == "C" {
		date = time.Date(int(paymentYear), time.Month(paymentMonth), (*account).DueDay, 0, 0, 0, 0, time.UTC)
	} else {
		date = time.Now()
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
