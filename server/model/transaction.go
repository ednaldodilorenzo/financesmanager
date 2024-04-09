package model

import "time"

type Transaction struct {
	UserDependent
	IdCategory  uint64          `gorm:"type:int(10)" json:"categoryId"`
	Category    RelatedCategory `gorm:"ForeignKey:IdCategory;references:id" json:"category"`
	IdAccount   int             `gorm:"type:bigint(10)" json:"accountId"`
	Account     RelatedAccount  `gorm:"ForeignKey:IdAccount;references:id"  json:"account"`
	Description string          `gorm:"type:varchar(50)" json:"description"`
	Value       float32         `gorm:"type:float" json:"value"`
	Date        time.Time       `gorm:"type:date" json:"date"`
	IdInvoice   *string         `gorm:"type:int(10)" json:"invoiceId"`
}

func (Transaction) TableName() string {
	return "transaction"
}
