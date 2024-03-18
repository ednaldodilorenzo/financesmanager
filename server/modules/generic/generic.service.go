package generic

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type GenericService[V interface{}] interface {
	FindAll() ([]V, error)
	Create(*V) error
	Update(int, map[string]interface{}) error
	FindById(id int) (*V, error)
}

type GenericServiceStruct[V interface{}] struct {
	repository GenericRepository[V]
}

func NewGenericService[V interface{}]() GenericService[V] {
	return &GenericServiceStruct[V]{
		repository: NewGenericRepository[V](),
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

func (c *GenericServiceStruct[V]) Update(id int, data map[string]interface{}) error {

	currentItem, err := c.repository.FindById(id)

	if err != nil {
		return &util.ItemNotFoundError{Message: "Item not found!"}
	}

	if err := c.repository.Update(currentItem, data); err != nil {
		return err
	}

	return nil
}
