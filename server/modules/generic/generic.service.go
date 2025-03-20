package generic

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type GenericService[V model.IUserDependent] interface {
	FindAll(userId int) ([]V, error)
	FindAllPaginatedAndFiltered(userId, limit, offset int, filter string) (*PaginatedResponse[V], error)
	Create(*V) error
	CreateAll([]V) error
	Update(id int, item *V, userId int) error
	FindById(id, userId int) (*V, error)
	DeleteRecord(id, userId int) error
}

type genericService[V model.IUserDependent] struct {
	repository GenericRepository[V]
}

func NewGenericService[V model.IUserDependent](repository GenericRepository[V]) GenericService[V] {
	return &genericService[V]{
		repository,
	}
}

func (c *genericService[V]) FindAll(userId int) ([]V, error) {
	return c.repository.FindAll(userId)
}

func (c *genericService[V]) FindAllPaginatedAndFiltered(userId, limit, offset int, filter string) (*PaginatedResponse[V], error) {
	return c.repository.FindAllPaginatedAndFiltered(userId, limit, offset, filter)
}

func (c *genericService[V]) CreateAll(items []V) error {
	return c.repository.CreateAll(items)
}

func (c *genericService[V]) DeleteRecord(id, userId int) error {
	return c.repository.Delete(id, userId)
}

func (c *genericService[V]) FindById(id, userId int) (*V, error) {
	currentItem, err := c.repository.FindById(id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, &util.ItemNotFoundError{Message: "Item not found!"}
		} else {
			return nil, err
		}
	}

	return currentItem, nil
}

func (c *genericService[V]) Create(item *V) error {
	return c.repository.Transaction(func(repo GenericRepository[V]) error {
		return repo.Create(item)
	})
}

func (c *genericService[V]) Update(id int, item *V, userId int) error {

	if err := c.repository.Update(id, item, userId); err != nil {
		return err
	}

	return nil
}
