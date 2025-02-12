package account

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type AccountService interface {
	generic.GenericService[*model.Account]
	FindByName(name string, userId int) (*model.Account, error)
}

type AccountServiceStruct struct {
	generic.GenericService[*model.Account]
	repository AccountRepository
}

func NewAccountService(service generic.GenericService[*model.Account], repository AccountRepository) AccountService {
	return &AccountServiceStruct{
		service,
		repository,
	}
}

func (as *AccountServiceStruct) FindByName(name string, userId int) (*model.Account, error) {
	return as.repository.FindByName(name, userId)
}
