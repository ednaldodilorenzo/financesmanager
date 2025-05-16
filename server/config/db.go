package config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connect(settings *DBSettings) error {
	var err error

	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", settings.Host, settings.Username, settings.Password, settings.DBName, settings.Port)

	d.DB, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

type Tx interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}

type GormTx struct {
	Tx *gorm.DB
}

func (t *GormTx) Commit() error {
	return t.Tx.Commit().Error
}

func (t *GormTx) Rollback() error {
	return t.Tx.Rollback().Error
}

type GormTxManager struct {
	dbConfig *Database
}

func (m *GormTxManager) Begin(ctx context.Context) (Tx, error) {
	tx := m.dbConfig.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTx{Tx: tx}, nil
}

func NewTxManager(dbConfig *Database) TxManager {
	return &GormTxManager{
		dbConfig: dbConfig,
	}
}
