package tag

import "github.com/ednaldo-dilorenzo/iappointment/model"

type TagService interface {
	FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error)
}

type tagService struct {
	TagRepository
}

func NewTagService(repository TagRepository) TagService {
	return &tagService{
		repository,
	}
}
