package model

import (
	"github.com/google/uuid"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"time"
)

func (m *TransactionDetail) TableName() string {
	return "transaction_detail"
}

type TransactionDetail struct {
	BaseModel

	Transaction   *Transaction `gorm:"foreignKey:TransactionID"`
	TransactionID string       `gorm:"column:transaction_id;size:36;"`
	AccountID     string       `gorm:"column:account_id;size:36;"`
	Type          string       `gorm:"column:type;size:100;"`
	Amount        int64        `gorm:"column:amount"`
	BalanceBefore int64        `gorm:"column:balance_before"`
	BalanceAfter  int64        `gorm:"column:balance_after"`
}

func (m *TransactionDetail) CreateNewTransactionDetail(transactionID, accountID, transactionType string, amount, balance int64) {
	m.ID = uuid.New().String()
	m.TransactionID = transactionID
	m.AccountID = accountID
	m.Type = transactionType
	m.Amount = amount
	m.BalanceBefore = balance
	if m.Type == constant.TransactionTypeCredit {
		m.BalanceAfter = m.BalanceBefore + amount
	} else {
		m.BalanceAfter = m.BalanceBefore - amount
	}
}

func (m *TransactionDetail) ToTransactionEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:            m.ID,
		Category:      m.Transaction.Category,
		Status:        m.Transaction.Status,
		Amount:        m.Transaction.FinalAmount,
		Type:          m.Type,
		BalanceBefore: m.BalanceBefore,
		BalanceAfter:  m.BalanceAfter,
		Description:   m.Transaction.Description,
		CreatedAt:     m.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Time.Format(time.RFC3339),
	}
}
