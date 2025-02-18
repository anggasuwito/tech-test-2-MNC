package usecase

import (
	"context"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/repository"
)

type UserAccUC interface {
	UpdateProfile(ctx context.Context, req *entity.UpdateProfileRequest) (*entity.UpdateProfileResponse, error)
}

type userAccUC struct {
	userAccRepo repository.UserAccountRepo
}

func NewUserAccUC(
	userAccRepo repository.UserAccountRepo,
) UserAccUC {
	return &userAccUC{
		userAccRepo: userAccRepo,
	}
}

func (u *userAccUC) UpdateProfile(ctx context.Context, req *entity.UpdateProfileRequest) (*entity.UpdateProfileResponse, error) {
	account, err := u.userAccRepo.GetUserAccountByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	account.UpdateUserAccount(req)
	err = u.userAccRepo.UpdateUserAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return &entity.UpdateProfileResponse{
		account.ToEntity(),
	}, nil
}
