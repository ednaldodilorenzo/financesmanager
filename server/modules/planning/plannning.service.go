package planning

import "github.com/ednaldo-dilorenzo/iappointment/model"

type PlanningService interface {
	FindByMonthAndYear(month int, year int, userId int) ([]model.Planning, error)
}

type planningService struct {
	repository PlanningRepository
}

func NewPlanningService(repository PlanningRepository) PlanningService {
	return &planningService{
		repository: repository,
	}
}

func (p *planningService) FindByMonthAndYear(month int, year int, userId int) ([]model.Planning, error) {
	return p.repository.FindByMonthAndYear(month, year, userId)
}
