package model

type Budget struct {
	UserDependent
	Year       uint  `gorm:"type:int" json:"year"`
	CategoryId uint8 `gorm:"type:int" json:"categoryId"`
	Value      int64 `gorm:"type:int(10)" json:"value"`
}

func (Budget) TableName() string {
	return "budget"
}
