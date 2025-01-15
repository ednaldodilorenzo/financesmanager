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
		Select("category.name, SUM(transaction.value) AS total").
		Joins("INNER JOIN category ON category.id = transaction.category_id").
		Where("EXTRACT(YEAR FROM transaction.payment_date) = ? AND EXTRACT(MONTH FROM transaction.payment_date) = ?", year, month).
		Group("category.name").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}

	if summaries == nil {
		summaries = []model.Planning{}
	}

	return summaries, nil
}
