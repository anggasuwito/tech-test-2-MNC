package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type TransactionRepo interface {
	GetTransactionByID(ctx context.Context, id string) (*model.Transaction, error)
	CreateTransactionDetail(ctx context.Context, data *model.TransactionDetail) error
	CreateTransaction(ctx context.Context, data *model.Transaction) error
	UpdateTransaction(ctx context.Context, req *model.Transaction) error
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

func (r *transactionRepo) CreateTransactionDetail(ctx context.Context, data *model.TransactionDetail) error {
	db := r.useTX(ctx)
	err := db.Debug().Create(data).Error
	if err != nil {
		return utils.ErrInternal("Failed create new transaction detail : "+err.Error(), "transactionRepo.CreateTransactionDetail")
	}
	return nil
}

func (r *transactionRepo) GetTransactionByID(ctx context.Context, id string) (*model.Transaction, error) {
	var data model.Transaction

	err := r.masterDB.
		Debug().
		Preload("TransactionDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc")
		}).
		Model(&model.Transaction{}).
		Where("deleted_at IS NULL").
		Where("id = ?", id).
		First(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Transaction Not Found", "transactionRepo.GetTransactionByID.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get transaction : "+err.Error(), "transactionRepo.GetTransactionByID")
	}

	return &data, nil
}

func (r *transactionRepo) UpdateTransaction(ctx context.Context, req *model.Transaction) error {
	db := r.useTX(ctx)
	err := db.Debug().Model(req).
		Clauses(clause.Returning{}).
		Where("id = ?", req.ID).
		Updates(req).
		Error
	if err != nil {
		return utils.ErrInternal("Failed update transaction : "+err.Error(), "userAccountRepo.UpdateTransaction")
	}
	return nil
}
