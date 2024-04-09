package config

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB
var DATABASE_URI string = "root:secret@tcp(localhost:3306)/finances?charset=utf8mb4&parseTime=True&loc=Local"

func Connect() error {
	var err error

	Database, err = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	return nil
}

type TransactionFunc func() error

func TxWrapper(fn TransactionFunc) error {
	var txHandle *gorm.DB

	if Database == nil {
		return errors.New("failed to initilize database")
	}

	txHandle = Database.Begin()

	defer func() {
		if r := recover(); r != nil {
			txHandle.Rollback()
		}
	}()

	err := fn()
	if err != nil {
		if e := txHandle.Rollback().Error; e != nil {
			return e
		}
		return err
	}
	if err = txHandle.Commit().Error; err != nil {
		return err
	}

	return nil
}
