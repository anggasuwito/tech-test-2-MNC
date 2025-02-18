package model

func (m *UserAccount) TableName() string {
	return "user_account"
}

type UserAccount struct {
	BaseModel

	PhoneNumber string `gorm:"column:phone_number;size:100;"`
	FirstName   string `gorm:"column:first_name;size:255;"`
	LastName    string `gorm:"column:last_name;size:255;"`
	Address     string `gorm:"column:address;size:255;"`
	Status      string `gorm:"column:status;size:100;"`
	PIN         string `gorm:"column:pin;size:255;"`
}
