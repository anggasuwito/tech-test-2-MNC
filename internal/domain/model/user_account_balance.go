package model

import (
	"github.com/google/uuid"
)

func (m *UserAccountBalance) TableName() string {
	return "user_account_balance"
}

type UserAccountBalance struct {
	BaseModel

	AccountID string `gorm:"column:account_id;size:36;"`
	Type      string `gorm:"column:type;size:100;"`
	Status    string `gorm:"column:status;size:100;"`
	Balance   int64  `gorm:"column:balance"`
}

func (m *UserAccountBalance) RegisterUserAccountBalance(userID string) {
	m.ID = uuid.New().String()
	m.AccountID = userID
	m.Type = "DEFAULT"
	m.Status = "ACTIVE"
	m.Balance = 0
}
