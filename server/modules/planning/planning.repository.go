package planning

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type PlanningRepository interface {
	FindByMonthAndYear(month int, year int) ([]model.Planning, error)
}

type PlanningRepositoryStruct struct {
	*config.Database
}

func NewPlanningRepository(database *config.Database) PlanningRepository {
	return &PlanningRepositoryStruct{
		database,
	}
}

func (p *PlanningRepositoryStruct) FindByMonthAndYear(month int, year int) ([]model.Planning, error) {
	var results []model.Planning

	subQuery := p.DB.Table("transaction t").
		Select("t.category_id, t.user_id, SUM(t.value) AS accumulated").
		Where("EXTRACT(YEAR FROM t.payment_date) = ?", year).
		Where("EXTRACT(MONTH FROM t.payment_date) <= ?", month).
		Group("t.category_id, t.user_id")

	err := p.DB.Table("transaction").
		Select(`
		    category.name AS name,
			category.type AS type,
			SUM(transaction.value) AS total, 
			budget.value AS planned, 
			a.accumulated
		`).
		Joins("INNER JOIN category ON category.id = transaction.category_id AND category.user_id = transaction.user_id").
		Joins("INNER JOIN (?) a ON a.category_id = transaction.category_id AND a.user_id = transaction.user_id", subQuery).
		Joins("LEFT JOIN budget ON budget.category_id = category.id AND category.user_id = budget.user_id").
		Where("EXTRACT(YEAR FROM transaction.payment_date) = ?", year).
		Where("EXTRACT(MONTH FROM transaction.payment_date) = ?", month).
		Where("budget.year = ?", year).
		Group("category.name, category.type, budget.value, a.accumulated").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []model.Planning{}
	}

	return results, nil
}
