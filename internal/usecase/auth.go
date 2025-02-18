package usecase

import (
	"context"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/repository"
	"tech-test-2-MNC/internal/utils"
)

type AuthUC interface {
	Login(ctx context.Context, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error)
	Register(ctx context.Context, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error)
}

type authUC struct {
	accountRepo repository.UserAccountRepo
}

func NewAuthUC(
	accountRepo repository.UserAccountRepo,
) AuthUC {
	return &authUC{
		accountRepo: accountRepo,
	}
}

func (u *authUC) Register(ctx context.Context, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error) {

	return nil, nil
}

func (u *authUC) Login(ctx context.Context, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error) {
	account, err := u.accountRepo.GetUserAccountByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	validPin := utils.CompareHashCredential(req.PIN, account.PIN)
	if !validPin {
		return nil, utils.ErrBadRequest("Invalid pin", "authUC.Login.CompareHashCredential")
	}

	accessToken, _, err := utils.GenerateJWT(account, constant.TokenTypeAccess)
	if err != nil {
		return nil, utils.ErrInternal("Failed generate jwt : "+err.Error(), "authUC.Login.GenerateJWT")
	}

	refreshToken, _, err := utils.GenerateJWT(account, constant.TokenTypeRefresh)
	if err != nil {
		return nil, utils.ErrInternal("Failed generate jwt : "+err.Error(), "authUC.Login.GenerateJWT")
	}

	return &entity.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
