package planning

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type PlanningRepository interface {
	FindByMonthAndYear(month int, year int) ([]model.Planning, error)
}

type PlanningRepositoryStruct struct {
	dbConfig *config.Database
}

func NewPlanningRepository(database *config.Database) PlanningRepository {
	return &PlanningRepositoryStruct{
		database,
	}
}

func (p *PlanningRepositoryStruct) FindByMonthAndYear(month int, year int) ([]model.Planning, error) {
	var summaries []model.Planning
	if err := p.dbConfig.DB.Model(&model.Transaction{}).
		Select("category.name, SUM(transaction.value) AS total, budget.value as planned").
		Joins("INNER JOIN category ON category.id = transaction.category_id AND category.user_id = transaction.user_id").
		Joins("LEFT JOIN budget ON budget.category_id = category.id AND category.user_id = budget.user_id").
		Where("EXTRACT(YEAR FROM transaction.payment_date) = ? AND EXTRACT(MONTH FROM transaction.payment_date) = ? AND budget.year = ?", year, month, year).
		Group("category.name, budget.value").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}

	if summaries == nil {
		summaries = []model.Planning{}
	}

	return summaries, nil
}
