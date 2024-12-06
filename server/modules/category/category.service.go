package category

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type CategoryService interface {
	generic.GenericService[*model.Category]
	FindByName(name string) (*model.Category, error)
}

type CategoryServiceStruct struct {
	generic.GenericService[*model.Category]
	repository CategoryRepository
}

func NewAccountService(service generic.GenericService[*model.Category], repository CategoryRepository) CategoryService {
	return &CategoryServiceStruct{
		service,
		repository,
	}
}

func (cs *CategoryServiceStruct) FindByName(name string) (*model.Category, error) {
	return cs.repository.FindByName(name)
}
