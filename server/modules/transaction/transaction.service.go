package transaction

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type TransactionService interface {
	generic.GenericService[model.Transaction]
}
