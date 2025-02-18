package model

import (
	"github.com/google/uuid"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/utils"
)

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

func (m *UserAccount) RegisterUserAccount(req *entity.AuthRegisterRequest, hashedPIN string) {
	m.ID = uuid.New().String()
	m.PhoneNumber = req.PhoneNumber
	m.FirstName = req.FirstName
	m.LastName = req.LastName
	m.Address = req.Address
	m.Status = "ACTIVE"
	m.PIN = hashedPIN
}

func (m *UserAccount) ToEntity() *entity.UserAccount {
	return &entity.UserAccount{
		ID:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		PhoneNumber: m.PhoneNumber,
		Address:     m.Address,
		CreatedAt:   utils.ParseTime(m.CreatedAt.Time),
		UpdatedAt:   utils.ParseTime(m.UpdatedAt.Time),
	}
}
