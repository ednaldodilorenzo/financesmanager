package generic

import (
	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"gorm.io/gorm"
)

type GenericRepository[V model.IUserDependent] interface {
	FindAll() ([]V, error)
	Create(*V) error
	Update(id int, item *V) error
	FindById(id int) (*V, error)
	Transaction(fn func(repo GenericRepository[V]) error) error //https://gist.github.com/IamNator/f1e9e6b1ae4d9e3eb66c73998f545f6c
}

type GenericRepositoryStruct[V model.IUserDependent] struct {
	dbConfig *config.Database
}

func NewGenericRepository[V model.IUserDependent](database *config.Database) GenericRepository[V] {
	return &GenericRepositoryStruct[V]{
		dbConfig: database,
	}
}

func (g *GenericRepositoryStruct[V]) FindById(id int) (*V, error) {
	var item V

	if err := g.dbConfig.DB.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (g *GenericRepositoryStruct[V]) FindAll() ([]V, error) {
	var items []V

	if err := g.dbConfig.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (g *GenericRepositoryStruct[V]) Create(item *V) error {
	if err := g.dbConfig.DB.Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (g *GenericRepositoryStruct[V]) Update(id int, item *V) error {
	if err := g.dbConfig.DB.Model(&item).Where("id = ?", id).Updates(item).Error; err != nil {
		return err
	}

	return nil
}

func (c *GenericRepositoryStruct[V]) withTx(tx *gorm.DB) GenericRepository[V] {
	return &GenericRepositoryStruct[V]{
		dbConfig: &config.Database{
			DB: tx,
		},
	}
}

func (c *GenericRepositoryStruct[V]) Transaction(fn func(repo GenericRepository[V]) error) error {
	tx := c.dbConfig.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	repo := c.withTx(tx)
	err := fn(repo)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
