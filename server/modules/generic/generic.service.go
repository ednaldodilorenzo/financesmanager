package generic

import (
	"context"
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
)

type GenericService[V model.IUserDependent] interface {
	FindAll(context.Context, int) ([]V, error)
	FindAllPaginatedAndFiltered(ctx context.Context, userId, limit, offset int, filter string) (*PaginatedResponse[V], error)
	Create(context.Context, *V) error
	CreateAll(context.Context, []V) error
	Update(context.Context, int, *V, int) error
	FindById(ctx context.Context, id, userId int) (*V, error)
	DeleteRecord(ctx context.Context, id, userId int) error
}

type genericService[V model.IUserDependent] struct {
	repository GenericRepository[V]
	txManager  config.TxManager
}

func NewGenericService[V model.IUserDependent](repository GenericRepository[V], txManager config.TxManager) GenericService[V] {
	return &genericService[V]{
		repository,
		txManager,
	}
}

func (c *genericService[V]) FindAll(ctx context.Context, userId int) ([]V, error) {
	return c.repository.FindAll(ctx, userId)
}

func (c *genericService[V]) FindAllPaginatedAndFiltered(ctx context.Context, userId, limit, offset int, filter string) (*PaginatedResponse[V], error) {
	return c.repository.FindAllPaginatedAndFiltered(ctx, userId, limit, offset, filter)
}

func (c *genericService[V]) CreateAll(ctx context.Context, items []V) error {
	tx, err := c.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	return c.repository.CreateAll(ctx, gormTx, items)
}

func (c *genericService[V]) DeleteRecord(ctx context.Context, id, userId int) error {
	tx, err := c.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = c.repository.Delete(ctx, gormTx, id, userId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *genericService[V]) FindById(ctx context.Context, id, userId int) (*V, error) {
	currentItem, err := c.repository.FindById(ctx, id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, &util.ItemNotFoundError{Message: "Item not found!"}
		} else {
			return nil, err
		}
	}

	return currentItem, nil
}

func (c *genericService[V]) Create(ctx context.Context, item *V) error {
	tx, err := c.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = c.repository.Create(ctx, gormTx, item)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *genericService[V]) Update(ctx context.Context, id int, item *V, userId int) error {
	tx, err := c.txManager.Begin(ctx)
	if err != nil {
		return err
	}

	gormTx := tx.(*config.GormTx).Tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = c.repository.Update(ctx, gormTx, id, item, userId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
