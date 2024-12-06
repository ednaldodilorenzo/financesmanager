package model

type Category struct {
	UserDependent
	Name string `gorm:"type:varchar(45);not null" json:"name"`
	Type string `gorm:"type:varchar(1);not null"  json:"type"`
}

type RelatedCategory struct {
	ID   uint64
	Name string `json:"name"`
	Type string `json:"type"`
}

func (Category) TableName() string {
	return "category"
}

func (RelatedCategory) TableName() string {
	return "category"
}
