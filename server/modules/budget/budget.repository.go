package budget

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
)

type BudgetRepository interface {
	generic.GenericRepository[*model.Budget]
	FindAllByYear(year, userId int) ([]model.Budget, error)
}

type BudgetRepositoryStruct struct {
	generic.GenericRepository[*model.Budget]
	*config.Database
}

func NewBudgetRepository(repository generic.GenericRepository[*model.Budget], database *config.Database) BudgetRepository {
	return &BudgetRepositoryStruct{
		repository,
		database,
	}
}

func (b *BudgetRepositoryStruct) FindAllByYear(year, userId int) ([]model.Budget, error) {
	var items []model.Budget

	if err := b.DB.Model(&items).Where("year = ? AND user_id = ?", year, userId).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
