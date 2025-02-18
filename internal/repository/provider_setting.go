package repository

import (
	"context"
	"gorm.io/gorm"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/utils"
)

type ProviderSettingRepo interface {
	GetProviderSetting(ctx context.Context, providerID string, settingType string) (*model.ProviderSetting, error)
}

type providerSettingRepo struct {
	masterDB *gorm.DB
}

func NewProviderSettingRepo(masterDB *gorm.DB) ProviderSettingRepo {
	return &providerSettingRepo{
		masterDB: masterDB,
	}
}

func (r *providerSettingRepo) GetProviderSetting(ctx context.Context, providerID string, settingType string) (*model.ProviderSetting, error) {
	var data model.ProviderSetting

	err := r.masterDB.
		Debug().
		Model(&model.ProviderSetting{}).
		Where("deleted_at IS NULL").
		Where("provider_id = ?", providerID).
		Where("category = ?", settingType).
		First(&data).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound("Provider Setting Not Found", "providerSettingRepo.GetProviderSetting.ErrRecordNotFound")
		}
		return nil, utils.ErrInternal("Failed get provider setting : "+err.Error(), "providerSettingRepo.GetProviderSetting")
	}

	return &data, nil
}
