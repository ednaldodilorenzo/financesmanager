package model

type Budget struct {
	UserDependent
	Year       uint  `gorm:"primay_key" json:"year"`
	CategoryId uint8 `gorm:"primay_key" json:"categoryId"`
	Value      int32 `gorm:"type:int(10)" json:"value"`
}

func (Budget) TableName() string {
	return "budget"
}
