package model

import "gorm.io/gorm"

type Category struct {
	UserDependent
	Name   string `gorm:"type:varchar(45);not null" json:"name"`
	Type   string `gorm:"type:char(1);not null"  json:"type"`
	Filter string `gorm:"type:varchar(100);not null"`
}

type RelatedCategory struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (Category) TableName() string {
	return "category"
}

func (RelatedCategory) TableName() string {
	return "category"
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.Filter = c.Name
	return
}

func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	c.Filter = c.Name
	return
}
