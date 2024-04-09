package generic

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type GenericRepository[V model.IUserDependent] interface {
	FindAll() ([]V, error)
	Create(*V) error
	Update(id int, item *V) error
	FindById(id int) (*V, error)
}

type GenericRepositoryStruct[V model.IUserDependent] struct {
}

func NewGenericRepository[V model.IUserDependent]() GenericRepository[V] {
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

func (c *GenericRepositoryStruct[V]) Update(id int, item *V) error {
	if err := config.Database.Model(&item).Where("id = ?", id).Updates(item).Error; err != nil {
		return err
	}

	return nil
}
