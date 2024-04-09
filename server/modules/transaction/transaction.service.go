package transaction

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type TransactionService interface {
	generic.GenericService[*model.Transaction]
	FindAllRelated(*int, *int) ([]model.Transaction, error)
}

type TransactionServiceStruct struct {
	generic.GenericService[*model.Transaction]
	repository TransactionRepository
}

func NewTransactionService(service generic.GenericService[*model.Transaction], repository TransactionRepository) TransactionService {
	return &TransactionServiceStruct{
		service,
		repository,
	}
}

func (ts *TransactionServiceStruct) FindAllRelated(month *int, year *int) ([]model.Transaction, error) {
	return ts.repository.FindAllWithRelationships(month, year)
}
