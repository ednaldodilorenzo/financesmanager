package model

type Account struct {
	UserDependent
	Name   string `gorm:"type:varchar(45);not null" json:"name"`
	Type   string `gorm:"type:char(1);not null" json:"type"`
	DueDay int    `gorm:"type:int" json:"dueDay"`
}

type RelatedAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (Account) TableName() string {
	return "account"
}

func (RelatedAccount) TableName() string {
	return "account"
}
