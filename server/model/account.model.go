package model

type Account struct {
	ID     int    `gorm:"autoIncrement;primary_key" json:"id"`
	IdUser uint64 `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(45);not null" json:"name"`
	Type   string `gorm:"type:char(1);not null" json:"type"`
	DueDay int    `gorm:"type:int" json:"dueDay"`
}

func (Account) TableName() string {
	return "account"
}
