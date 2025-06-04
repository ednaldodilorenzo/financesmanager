package generic

import (
	"context"
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaginatedResponse[V model.IUserDependent] struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Items []V `json:"items"`
}

type GenericRepository[V model.IUserDependent] interface {
	FindAll(context.Context, int) ([]V, error)
	FindAllPaginatedAndFiltered(ctx context.Context, userId, limit, offset int, filter string) (*PaginatedResponse[V], error)
	Create(context.Context, *gorm.DB, V) error
	CreateAll(context.Context, *gorm.DB, []V) error
	Update(context.Context, *gorm.DB, int, V, int) error
	FindById(context.Context, int, int) (V, error)
	Delete(context.Context, *gorm.DB, int, int) error
}

type genericRepository[V model.IUserDependent] struct {
	dbConfig *config.Database
}

func NewGenericRepository[V model.IUserDependent](database *config.Database, txManager config.TxManager) GenericRepository[V] {
	return &genericRepository[V]{
		dbConfig: database,
	}
}

func (g *genericRepository[V]) FindById(ctx context.Context, id int, userId int) (V, error) {
	var item V
	var zero V

	if err := g.dbConfig.DB.WithContext(ctx).First(&item, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		return zero, err
	}

	return item, nil
}

func (g *genericRepository[V]) Delete(ctx context.Context, db *gorm.DB, id, userId int) error {
	var item V

	if err := db.Delete(&item, id).Error; err != nil {
		return err
	}

	return nil
}

func (g *genericRepository[V]) FindAll(ctx context.Context, userId int) ([]V, error) {
	var items []V

	var zeroValue V

	query := g.dbConfig.DB.WithContext(ctx).Model(&zeroValue)

	query.Where("user_id = ?", userId)

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (g *genericRepository[V]) FindAllPaginatedAndFiltered(ctx context.Context, userId, limit, offset int, filter string) (*PaginatedResponse[V], error) {
	var totalCount int64
	var items []V
	var zeroValue V

	query := g.dbConfig.DB.WithContext(ctx).Model(&zeroValue)
	query = query.Where("user_id = ?", userId)

	// Apply filter if it's not empty
	if filter != "" {
		filterString := "%" + strings.ToLower(filter) + "%"
		query = query.Where("LOWER(filter) LIKE ?", filterString) // Replace "name" with the actual column you want to filter by
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

func (g *genericRepository[V]) Create(ctx context.Context, db *gorm.DB, item V) error {
	if err := db.Clauses(clause.Returning{}).Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (g *genericRepository[V]) CreateAll(ctx context.Context, db *gorm.DB, items []V) error {
	if err := db.Create(&items).Error; err != nil {
		return err
	}

	return nil
}

func (g *genericRepository[V]) Update(ctx context.Context, db *gorm.DB, id int, item V, userId int) error {
	if err := db.Model(&item).Where("id = ? AND user_id = ?", id, userId).Updates(item).Error; err != nil {
		return err
	}

	return nil
}
