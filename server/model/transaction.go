package model

import "time"

type Transaction struct {
	ID          uint32    `gorm:"autoIncrement;primary_key" json:"id"`
	IdUser      uint64    `gorm:"type:bigint(10);primary_key" json:"userId"`
	IdCategory  uint64    `gorm:"type:int(10)" json:"categoryId"`
	Category    Category  `gorm:"ForeignKey:IdCategory;references:id"`
	IdAccount   int       `gorm:"type:bigint(10)" json:"accountId"`
	Account     Account   `gorm:"ForeignKey:IdAccount;references:id"`
	Description string    `gorm:"type:varchar(50)" json:"description"`
	Value       float32   `gorm:"type:float" json:"value"`
	Date        time.Time `gorm:"type:date" json:"date"`
	IdInvoice   *string   `gorm:"type:int(10)" json:"invoiceId"`
}

func (Transaction) TableName() string {
	return "transaction"
}
