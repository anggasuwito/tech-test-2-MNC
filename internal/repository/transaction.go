package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type TransactionRepo interface {
	GetTransactionByID(ctx context.Context, id string) (*model.Transaction, error)
	GetTransactionDetails(ctx context.Context, filter []*entity.Filter, sort []*entity.Filter, page, limit int64) ([]*model.TransactionDetail, int64, error)
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

func (r *transactionRepo) GetTransactionDetails(ctx context.Context, filter []*entity.Filter, sort []*entity.Filter, page, limit int64) ([]*model.TransactionDetail, int64, error) {
	var (
		res   = []*model.TransactionDetail{}
		count int64
		err   error
	)

	q := r.masterDB.Debug().Preload("Transaction").Model(&model.TransactionDetail{}).Where("deleted_at is NULL")
	for _, v := range filter {
		if v.Value == "" {
			continue
		}

		valSlice := strings.Split(v.Value, "|")

		switch v.Field {
		case "account_id":
			if len(valSlice) > 1 {
				q.Where(fmt.Sprintf("%s IN(?)", v.Field), valSlice)
			} else {
				q.Where(fmt.Sprintf("%s = ?", v.Field), v.Value)
			}
		}
	}

	for _, s := range sort {
		q.Order(fmt.Sprintf("%s %s", s.Field, s.Value))
	}

	if len(sort) < 1 {
		q.Order("created_at DESC")
	}

	q.Count(&count)

	err = q.Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&res).Error
	return res, count, err
}
