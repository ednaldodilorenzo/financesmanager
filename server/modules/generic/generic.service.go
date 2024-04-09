package generic

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type GenericService[V model.IUserDependent] interface {
	FindAll() ([]V, error)
	Create(*V) error
	Update(id int, item *V) error
	FindById(id int) (*V, error)
}

type GenericServiceStruct[V model.IUserDependent] struct {
	repository GenericRepository[V]
}

func NewGenericService[V model.IUserDependent](repository GenericRepository[V]) GenericService[V] {
	return &GenericServiceStruct[V]{
		repository: repository,
	}
}

func (c *GenericServiceStruct[V]) FindAll() ([]V, error) {
	return c.repository.FindAll()
}

func (c *GenericServiceStruct[V]) FindById(id int) (*V, error) {
	currentItem, err := c.repository.FindById(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, &util.ItemNotFoundError{Message: "Item not found!"}
		} else {
			return nil, err
		}
	}

	return currentItem, nil
}

func (c *GenericServiceStruct[V]) Create(item *V) error {
	return c.repository.Create(item)
}

func (c *GenericServiceStruct[V]) Update(id int, item *V) error {

	if err := c.repository.Update(id, item); err != nil {
		return err
	}

	return nil
}
