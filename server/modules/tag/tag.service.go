package tag

import "github.com/ednaldo-dilorenzo/iappointment/model"

type TagService interface {
	FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error)
}

type tagService struct {
	repository TagRepository
}

func NewTagService(repository TagRepository) TagService {
	return &tagService{
		repository,
	}
}

func (t *tagService) FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error) {
	return t.repository.FindAllWithFilter(filter, userId)
}
