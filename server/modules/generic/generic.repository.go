package generic

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
)

type GenericRepository[V interface{}] interface {
	FindAll() ([]V, error)
	Create(*V) error
	Update(*V, map[string]interface{}) error
	FindById(id int) (*V, error)
}

type GenericRepositoryStruct[V interface{}] struct {
}

func NewGenericRepository[V interface{}]() GenericRepository[V] {
	return &GenericRepositoryStruct[V]{}
}

func (g *GenericRepositoryStruct[V]) FindById(id int) (*V, error) {
	var item V

	if err := config.Database.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *GenericRepositoryStruct[V]) FindAll() ([]V, error) {
	var items []V

	if err := config.Database.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (c *GenericRepositoryStruct[V]) Create(item *V) error {
	if err := config.Database.Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (c *GenericRepositoryStruct[V]) Update(item *V, data map[string]interface{}) error {
	if err := config.Database.Model(&item).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
