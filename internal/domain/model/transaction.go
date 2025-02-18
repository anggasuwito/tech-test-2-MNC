package model

func (m *Transaction) TableName() string {
	return "transaction"
}

type Transaction struct {
	BaseModel

	Category    string `gorm:"column:category;size:100;"`
	Status      string `gorm:"column:status;size:100;"`
	Description string `gorm:"column:description;size:255;"`
	Amount      int64  `gorm:"column:amount"`
	FinalAmount int64  `gorm:"column:final_amount"`
}
