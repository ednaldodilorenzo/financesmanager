package budget

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type BudgetService interface {
	generic.GenericService[*model.Budget]
	FindAllByYear(year, userId int) ([]model.Budget, error)
}

type BudgetServiceStruct struct {
	generic.GenericService[*model.Budget]
	repository BudgetRepository
}

func NewBudgetService(service generic.GenericService[*model.Budget], repository BudgetRepository) BudgetService {
	return &BudgetServiceStruct{
		service,
		repository,
	}
}

func (b *BudgetServiceStruct) FindAllByYear(year, userId int) ([]model.Budget, error) {
	return b.repository.FindAllByYear(year, userId)
}
