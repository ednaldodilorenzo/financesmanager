package budget

import (
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type BudgetService interface {
	generic.GenericService[*model.Budget]
	FindAllByYear(year, userId int) ([]model.Budget, error)
}

type budgetService struct {
	generic.GenericService[*model.Budget]
	repository BudgetRepository
}

func NewBudgetService(service generic.GenericService[*model.Budget], repository BudgetRepository) BudgetService {
	return &budgetService{
		service,
		repository,
	}
}

func (b *budgetService) FindAllByYear(year, userId int) ([]model.Budget, error) {
	return b.repository.FindAllByYear(year, userId)
}
