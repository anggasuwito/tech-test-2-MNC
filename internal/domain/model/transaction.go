package model

import "database/sql"

func (m *Transaction) TableName() string {
	return "transaction"
}

type Transaction struct {
	BaseModel

	TransactionDetails []*TransactionDetail `gorm:"foreignKey:ID"`
	Category           string               `gorm:"column:category;size:100;"`
	Status             string               `gorm:"column:status;size:100;"`
	Description        string               `gorm:"column:description;size:255;"`
	Amount             int64                `gorm:"column:amount"`
	FinalAmount        int64                `gorm:"column:final_amount"`
	FailedReason       string               `gorm:"column:failed_reason;size:255;"`
	CompletedAt        sql.NullTime         `gorm:"column:completed_at"`
}
