package model

import (
	"database/sql"
	"github.com/google/uuid"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"time"
)

func (m *Transaction) TableName() string {
	return "transaction"
}

type Transaction struct {
	BaseModel

	TransactionDetails []*TransactionDetail `gorm:"foreignKey:TransactionID"`
	Category           string               `gorm:"column:category;size:100;"`
	Status             string               `gorm:"column:status;size:100;"`
	Description        string               `gorm:"column:description;size:255;"`
	Amount             int64                `gorm:"column:amount"`
	FinalAmount        int64                `gorm:"column:final_amount"`
	FailedReason       string               `gorm:"column:failed_reason;size:255;"`
	CompletedAt        sql.NullTime         `gorm:"column:completed_at"`
}

func (m *Transaction) ToEntity(transactionType string, balanceBefore, balanceAfter int64) *entity.Transaction {
	return &entity.Transaction{
		ID:            m.ID,
		Category:      m.Category,
		Status:        m.Status,
		Amount:        m.FinalAmount,
		Type:          transactionType,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   m.Description,
		CreatedAt:     m.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func (m *Transaction) CreateNewTransaction(category, description string, amount int64) {
	m.ID = uuid.New().String()
	m.Category = category
	m.Status = constant.TransactionStatusPending
	m.Description = description
	m.Amount = amount
	m.FinalAmount = amount
}
