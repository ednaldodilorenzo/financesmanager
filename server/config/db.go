package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
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

	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", settings.Username, settings.Password, settings.Host, settings.Port, settings.DBName)

	d.DB, err = gorm.Open(mysql.Open(dbUri), &gorm.Config{
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
