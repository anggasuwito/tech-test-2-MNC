package repository

import (
	"context"
	"gorm.io/gorm"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type UserAccountRepo interface {
	GetUserAccountByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
	GetUserAccountByID(ctx context.Context, id string) (*model.UserAccount, error)
	GetUserAccountByVA(ctx context.Context, vaNumber string) (*model.UserAccount, error)
	UpdateUserBalance(ctx context.Context, id string, amount int64, updateType string) error
}

type userAccountRepo struct {
	masterDB *gorm.DB
}

func NewUserAccountRepo(masterDB *gorm.DB) UserAccountRepo {
	return &userAccountRepo{
		masterDB: masterDB,
	}
}

func (r *userAccountRepo) useTX(ctx context.Context) *gorm.DB {
	if tx := utils.GetTransactionFromContext(ctx); tx != nil {
		return tx
	}
	return r.masterDB
}

func (r *userAccountRepo) GetUserAccountByPhone(ctx context.Context, phone string) (*model.UserAccount, error) {
	var data model.UserAccount

	err := r.masterDB.
		Debug().
		Model(&model.UserAccount{}).
		Where("deleted_at IS NULL").
		Where("phone = ?", phone).
		First(&data).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound("User Account Not Found", "userAccountRepo.GetUserAccountByPhone.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get user account : "+err.Error(), "userAccountRepo.GetUserAccountByPhone")
	}

	return &data, nil
}

func (r *userAccountRepo) GetUserAccountByID(ctx context.Context, id string) (*model.UserAccount, error) {
	var data model.UserAccount

	err := r.masterDB.
		Debug().
		Model(&model.UserAccount{}).
		Where("deleted_at IS NULL").
		Where("id = ?", id).
		First(&data).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound("User Account Not Found", "userAccountRepo.GetUserAccountByID.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get user account : "+err.Error(), "userAccountRepo.GetUserAccountByID")
	}

	return &data, nil
}

func (r *userAccountRepo) GetUserAccountByVA(ctx context.Context, vaNumber string) (*model.UserAccount, error) {
	var data model.UserAccount

	err := r.masterDB.
		Debug().
		Model(&model.UserAccount{}).
		Where("deleted_at IS NULL").
		Where("va_number = ?", vaNumber).
		First(&data).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound("User Account Not Found", "userAccountRepo.GetUserAccountByVA.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get user account : "+err.Error(), "userAccountRepo.GetUserAccountByVA")
	}

	return &data, nil
}

func (r *userAccountRepo) UpdateUserBalance(ctx context.Context, id string, amount int64, updateType string) error {
	db := r.useTX(ctx)
	q := db.
		Debug().
		Model(&model.UserAccount{}).
		Where("id", id)

	switch updateType {
	case constant.INCREASE:
		q.Update("balance", gorm.Expr("balance + ?", amount))
	case constant.DECREASE:
		q.Update("balance", gorm.Expr("balance - ?", amount))
	default:
		return utils.ErrBadRequest("Invalid update balance type", "userAccountRepo.UpdateUserBalance.Type")
	}

	err := q.Error
	if err != nil {
		return utils.ErrInternal("Failed update balance : "+err.Error(), "userAccountRepo.UpdateUserBalance.Update")
	}

	return nil
}
