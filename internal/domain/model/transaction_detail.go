package model

func (m *TransactionDetail) TableName() string {
	return "transaction_detail"
}

type TransactionDetail struct {
	BaseModel

	TransactionID string `gorm:"column:transaction_id;size:36;"`
	AccountID     string `gorm:"column:account_id;size:36;"`
	Type          string `gorm:"column:type;size:100;"`
	Amount        int64  `gorm:"column:amount"`
	BalanceBefore int64  `gorm:"column:balance_before"`
	BalanceAfter  int64  `gorm:"column:balance_after"`
}
