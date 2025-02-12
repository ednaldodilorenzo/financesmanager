package model

import (
	"time"
)

type Tabler interface {
	TableName() string
}

type User struct {
	ID        uint64     `gorm:"autoIncrement;primary_key"`
	Name      string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
	Verified  bool       `gorm:"not null;default:false"`
}

func (User) TableName() string {
	return "user"
}
