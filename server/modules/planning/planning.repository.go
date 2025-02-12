package planning

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type PlanningRepository interface {
	FindByMonthAndYear(month int, year int, userId int) ([]model.Planning, error)
}

type PlanningRepositoryStruct struct {
	*config.Database
}

func NewPlanningRepository(database *config.Database) PlanningRepository {
	return &PlanningRepositoryStruct{
		database,
	}
}

func (p *PlanningRepositoryStruct) FindByMonthAndYear(month int, year int, userId int) ([]model.Planning, error) {
	var results []model.Planning

	accumulatedSubQuery := p.DB.Table("transaction t").
		Select("t.category_id, t.user_id, SUM(t.value) AS accumulated").
		Where("EXTRACT(YEAR FROM t.payment_date) = ?", year).
		Where("EXTRACT(MONTH FROM t.payment_date) <= ?", month).
		Group("t.category_id, t.user_id")

	monthYearSubQuery := p.DB.Table("transaction t").
		Select("t.category_id, t.user_id, SUM(t.value) AS value").
		Where("EXTRACT(YEAR FROM t.payment_date) = ?", year).
		Where("EXTRACT(MONTH FROM t.payment_date) = ?", month).
		Group("t.category_id, t.user_id")

	err := p.DB.Table("budget").
		Select(`
		    category.name AS name,
			category.type AS type,
			m.value AS total, 
			budget.value AS planned, 
			a.accumulated
		`).
		Joins("INNER JOIN category ON category.id = budget.category_id AND category.user_id = budget.user_id").
		Joins("LEFT JOIN (?) m ON m.category_id = budget.category_id AND m.user_id = budget.user_id", monthYearSubQuery).
		Joins("LEFT JOIN (?) a ON a.category_id = budget.category_id AND a.user_id = budget.user_id", accumulatedSubQuery).
		Where("budget.year = ? AND budget.user_id = ?", year, userId).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []model.Planning{}
	}

	return results, nil
}
