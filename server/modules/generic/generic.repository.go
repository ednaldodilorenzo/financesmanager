package generic

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"gorm.io/gorm"
)

type PaginatedResponse[V model.IUserDependent] struct {
	Page  int
	Total int
	Items []V
}

type GenericRepository[V model.IUserDependent] interface {
	FindAll() ([]V, error)
	FindAllPaginatedAndFiltered(limit, offset int, filter string) (*PaginatedResponse[V], error)
	Create(*V) error
	CreateAll([]V) error
	Update(id int, item *V) error
	FindById(id int) (*V, error)
	Delete(id int) error
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

func (g *GenericRepositoryStruct[V]) Delete(id int) error {
	var item V

	if err := g.dbConfig.DB.Delete(&item, id).Error; err != nil {
		return err
	}

	return nil
}

func (g *GenericRepositoryStruct[V]) FindAll() ([]V, error) {
	var items []V

	if err := g.dbConfig.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (g *GenericRepositoryStruct[V]) FindAllPaginatedAndFiltered(limit, offset int, filter string) (*PaginatedResponse[V], error) {
	var totalCount int64
	var items []V

	zeroValue := new(V)

	query := g.dbConfig.DB.Model(zeroValue)

	// Apply filter if it's not empty
	if filter != "" {
		query = query.Where("LOWER(filter) LIKE ?", "%"+strings.ToLower(filter)+"%") // Replace "name" with the actual column you want to filter by
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination and fetch results
	if err := query.Limit(limit).Offset(offset - 1).Find(&items).Error; err != nil {
		return nil, err
	}

	return &PaginatedResponse[V]{
		Page:  offset,
		Total: int(totalCount),
		Items: items,
	}, nil
}

func (g *GenericRepositoryStruct[V]) Create(item *V) error {
	if err := g.dbConfig.DB.Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (g *GenericRepositoryStruct[V]) CreateAll(items []V) error {
	if err := g.dbConfig.DB.Create(&items).Error; err != nil {
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
