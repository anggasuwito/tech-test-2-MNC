package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type UserAccountRepo interface {
	CheckExistPhone(ctx context.Context, phone string) (bool, error)
	GetUserAccountByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
	CreateUserAccount(ctx context.Context, req *model.UserAccount) error
	CreateUserAccountBalance(ctx context.Context, req *model.UserAccountBalance) error
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
		Where("phone_number = ?", phone).
		First(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("User Account Not Found", "userAccountRepo.GetUserAccountByPhone.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get user account : "+err.Error(), "userAccountRepo.GetUserAccountByPhone")
	}

	return &data, nil
}

func (r *userAccountRepo) CheckExistPhone(ctx context.Context, phone string) (bool, error) {
	var count int64

	err := r.masterDB.
		Debug().
		Model(&model.UserAccount{}).
		Where("deleted_at IS NULL").
		Where("phone_number = ?", phone).
		Count(&count).
		Error

	if err != nil {
		return false, utils.ErrInternal("Failed check exist phone : "+err.Error(), "userAccountRepo.CheckExistPhone")
	}

	return count > 0, nil
}

func (r *userAccountRepo) CreateUserAccount(ctx context.Context, req *model.UserAccount) error {
	db := r.useTX(ctx)
	err := db.Debug().Create(req).Error
	if err != nil {
		return utils.ErrInternal("Failed create new user account : "+err.Error(), "userAccountRepo.CreateUserAccount")
	}
	return nil
}

func (r *userAccountRepo) CreateUserAccountBalance(ctx context.Context, req *model.UserAccountBalance) error {
	db := r.useTX(ctx)
	err := db.Debug().Create(req).Error
	if err != nil {
		return utils.ErrInternal("Failed create new user account : "+err.Error(), "userAccountRepo.CreateUserAccountBalance")
	}
	return nil
}
