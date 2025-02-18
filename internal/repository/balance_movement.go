package repository

import (
	"context"
	"gorm.io/gorm"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type BalanceMovementRepo interface {
	GetUserAccBalanceMovementList(ctx context.Context, accountID string) ([]*model.BalanceMovement, error)
	CreateBalanceMovement(ctx context.Context, data *model.BalanceMovement) error
}

type balanceMovementRepo struct {
	masterDB *gorm.DB
}

func NewBalanceMovementRepo(masterDB *gorm.DB) BalanceMovementRepo {
	return &balanceMovementRepo{
		masterDB: masterDB,
	}
}

func (r *balanceMovementRepo) useTX(ctx context.Context) *gorm.DB {
	if tx := utils.GetTransactionFromContext(ctx); tx != nil {
		return tx
	}
	return r.masterDB
}

func (r *balanceMovementRepo) CreateBalanceMovement(ctx context.Context, data *model.BalanceMovement) error {
	db := r.useTX(ctx)
	err := db.Debug().Create(data).Error
	if err != nil {
		return utils.ErrInternal("Failed create new balance movement : "+err.Error(), "balanceMovementRepo.CreateBalanceMovement")
	}
	return nil
}

func (r *balanceMovementRepo) GetUserAccBalanceMovementList(ctx context.Context, accountID string) ([]*model.BalanceMovement, error) {
	var data []*model.BalanceMovement

	err := r.masterDB.
		Debug().
		Model(&model.BalanceMovement{}).
		Where("deleted_at IS NULL").
		Where("user_account_id = ?", accountID).
		Order("created_at DESC").
		Find(&data).
		Error
	if err != nil {
		return nil, utils.ErrInternal("Failed get account balance movement list account : "+err.Error(), "balanceMovementRepo.GetUserAccBalanceMovementList")
	}

	return data, nil
}
