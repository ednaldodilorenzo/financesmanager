package planning

import "github.com/ednaldo-dilorenzo/iappointment/model"

type PlanningService interface {
	FindByMonthAndYear(month int, year int) ([]model.Planning, error)
}

type PlanningServiceStruct struct {
	repository PlanningRepository
}

func NewPlanningService(repository PlanningRepository) PlanningService {
	return &PlanningServiceStruct{
		repository: repository,
	}
}

func (p *PlanningServiceStruct) FindByMonthAndYear(month int, year int) ([]model.Planning, error) {
	return p.repository.FindByMonthAndYear(month, year)
}
