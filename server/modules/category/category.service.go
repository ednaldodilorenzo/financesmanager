package category

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type CategoryService interface {
	generic.GenericService[*model.Category]
	FindByName(name string, userId int) (*model.Category, error)
}

type categoryService struct {
	generic.GenericService[*model.Category]
	repository CategoryRepository
}

func NewAccountService(service generic.GenericService[*model.Category], repository CategoryRepository) CategoryService {
	return &categoryService{
		service,
		repository,
	}
}

func (cs *categoryService) FindByName(name string, userId int) (*model.Category, error) {
	return cs.repository.FindByName(name, userId)
}
