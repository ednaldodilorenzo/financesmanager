package model

import "time"

type Transaction struct {
	UserDependent
	CategoryId      uint64           `gorm:"type:int(10)" json:"categoryId"`
	Category        RelatedCategory  `gorm:"ForeignKey:CategoryId;references:id" json:"category"`
	AccountId       int              `gorm:"type:bigint(10)" json:"accountId"`
	Account         RelatedAccount   `gorm:"ForeignKey:AccountId;references:id"  json:"account"`
	Description     string           `gorm:"type:varchar(50)" json:"description"`
	Value           int32            `gorm:"type:int(10)" json:"value"`
	PaymentDate     time.Time        `gorm:"type:date;column:payment_date" json:"paymentDate"` // Changed column name
	TransactionDate time.Time        `gorm:"type:date" json:"transactionDate"`
	Detail          *string          `gorm:"type:varchar(100)" json:"detail"`
	Tags            []TransactionTag `gorm:"foreignKey:TransactionId,UserId;references:Id,UserId;constraint:OnUpdate:CASCADE" json:"tags"`
}

func (Transaction) TableName() string {
	return "transaction"
}

type TransactionTag struct {
	TransactionId uint   `gorm:"primaryKey" json:"transactionId"`
	Tag           string `gorm:"primaryKey;size:255" json:"tag"`
	UserId        uint   `gorm:"primaryKey;" json:"userId"`
}

func (TransactionTag) TableName() string {
	return "transaction_tag"
}
