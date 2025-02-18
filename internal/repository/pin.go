package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"tech-test-2-MNC/internal/utils"
	"time"
)

type PINRepo interface {
	SetVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) error
	DeleteVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) error
	GetVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) (bool, error)
}

type pinRepo struct {
	redisClient *redis.Client
}

func NewPINRepo(redisClient *redis.Client) PINRepo {
	return &pinRepo{
		redisClient: redisClient,
	}
}

func (r *pinRepo) SetVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) error {
	key := utils.GetVerifiedPINKey(accountID, pinType)
	err := r.redisClient.Set(ctx, key, true, 5*time.Minute).Err()
	if err != nil {
		return utils.ErrInternal("Failed set verified pin cache : "+err.Error(), "pinRepo.SetVerifiedPINByTypeCache.Set")
	}
	return nil
}

func (r *pinRepo) GetVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) (bool, error) {
	key := utils.GetVerifiedPINKey(accountID, pinType)
	exist, err := r.redisClient.Get(ctx, key).Bool()
	if err != nil && err != redis.Nil {
		return false, utils.ErrInternal("Failed get verified pin cache : "+err.Error(), "pinRepo.GetVerifiedPINByTypeCache.Get")
	}
	return exist, nil
}

func (r *pinRepo) DeleteVerifiedPINByTypeCache(ctx context.Context, accountID string, pinType string) error {
	key := utils.GetVerifiedPINKey(accountID, pinType)
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return utils.ErrInternal("Failed delete verified pin cache : "+err.Error(), "pinRepo.DeleteVerifiedPINByTypeCache.Del")
	}
	return nil
}
