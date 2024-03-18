package model

type Category struct {
	ID     uint64 `gorm:"autoIncrement;primary_key" json:"id"`
	Name   string `gorm:"type:varchar(45);not null" json:"name"`
	Type   string `gorm:"type:varchar(1);not null"  json:"type"`
	IdUser uint64
}

func (Category) TableName() string {
	return "category"
}
