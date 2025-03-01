package tag

import "github.com/ednaldo-dilorenzo/iappointment/model"

type TagService interface {
	FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error)
}

type TagServiceStruct struct {
	repository TagRepository
}

func NewTagService(repository TagRepository) TagService {
	return &TagServiceStruct{
		repository,
	}
}

func (t *TagServiceStruct) FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error) {
	return t.repository.FindAllWithFilter(filter, userId)
}
