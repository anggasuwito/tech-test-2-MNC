package repository

import (
	"context"
	"gorm.io/gorm"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type TransactionRepo interface {
	CreateTransaction(ctx context.Context, data *model.Transaction) error
}

type transactionRepo struct {
	masterDB *gorm.DB
}

func NewTransactionRepo(masterDB *gorm.DB) TransactionRepo {
	return &transactionRepo{
		masterDB: masterDB,
	}
}

func (r *transactionRepo) useTX(ctx context.Context) *gorm.DB {
	if tx := utils.GetTransactionFromContext(ctx); tx != nil {
		return tx
	}
	return r.masterDB
}

func (r *transactionRepo) CreateTransaction(ctx context.Context, data *model.Transaction) error {
	db := r.useTX(ctx)
	err := db.Debug().Create(data).Error
	if err != nil {
		return utils.ErrInternal("Failed create new transaction : "+err.Error(), "transactionRepo.CreateTransaction")
	}
	return nil
}
