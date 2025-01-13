package model

import "time"

type Transaction struct {
	UserDependent
	IdCategory      uint64           `gorm:"type:int(10)" json:"categoryId"`
	Category        RelatedCategory  `gorm:"ForeignKey:IdCategory;references:id" json:"category"`
	IdAccount       int              `gorm:"type:bigint(10)" json:"accountId"`
	Account         RelatedAccount   `gorm:"ForeignKey:IdAccount;references:id"  json:"account"`
	Description     string           `gorm:"type:varchar(50)" json:"description"`
	Value           int32            `gorm:"type:int(10)" json:"value"`
	PaymentDate     time.Time        `gorm:"type:date;column:payment_date" json:"paymentDate"` // Changed column name
	TransactionDate time.Time        `gorm:"type:date" json:"transactionDate"`
	Detail          *string          `gorm:"type:varchar(100)" json:"detail"`
	Tags            []TransactionTag `gorm:"foreignKey:TransactionID,UserID;references:ID,IdUser" json:"tags"`
}

func (Transaction) TableName() string {
	return "transaction"
}

type TransactionTag struct {
	TransactionID uint   `gorm:"primaryKey" json:"transactionId"`
	Tag           string `gorm:"primaryKey;size:255" json:"tag"`
	UserID        uint   `gorm:"primaryKey" json:"userId"`
}

func (TransactionTag) TableName() string {
	return "transaction_tag"
}
