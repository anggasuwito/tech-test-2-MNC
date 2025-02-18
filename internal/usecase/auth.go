package usecase

import (
	"context"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/repository"
	"tech-test-2-MNC/internal/utils"
)

type AuthUC interface {
	Login(ctx context.Context, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error)
	Register(ctx context.Context, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error)
}

type authUC struct {
	txWrapper   repository.TransactionWrapper
	accountRepo repository.UserAccountRepo
}

func NewAuthUC(
	txWrapper repository.TransactionWrapper,
	accountRepo repository.UserAccountRepo,
) AuthUC {
	return &authUC{
		txWrapper:   txWrapper,
		accountRepo: accountRepo,
	}
}

func (u *authUC) Register(ctx context.Context, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error) {
	var (
		userAccount        = &model.UserAccount{}
		userAccountBalance = &model.UserAccountBalance{}
	)

	exist, err := u.accountRepo.CheckExistPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, utils.ErrBadRequest("Phone number already registered", "authUC.Register.exist")
	}

	hashedPIN, err := utils.HashCredential(req.PIN)
	if err != nil {
		return nil, utils.ErrInternal("Failed hash credential : "+err.Error(), "authUC.Register.utils.HashCredential")
	}

	if err = u.txWrapper.ExecuteTransaction(ctx,
		func(ctxTX context.Context) error {
			//create user account
			userAccount.RegisterUserAccount(req, hashedPIN)
			err = u.accountRepo.CreateUserAccount(ctxTX, userAccount)
			if err != nil {
				return err
			}

			//create user account balance
			userAccountBalance.RegisterUserAccountBalance(userAccount.ID)
			err = u.accountRepo.CreateUserAccountBalance(ctxTX, userAccountBalance)
			if err != nil {
				return err
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return &entity.AuthRegisterResponse{
		userAccount.ToEntity(),
	}, nil
}

func (u *authUC) Login(ctx context.Context, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error) {
	account, err := u.accountRepo.GetUserAccountByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	validPin := utils.CompareHashCredential(req.PIN, account.PIN)
	if !validPin {
		return nil, utils.ErrBadRequest("Invalid pin", "authUC.Login.CompareHashCredential")
	}

	accessToken, _, err := utils.GenerateJWT(account.ToJWTAccInfo(), constant.TokenTypeAccess)
	if err != nil {
		return nil, utils.ErrInternal("Failed generate jwt : "+err.Error(), "authUC.Login.GenerateJWT")
	}

	refreshToken, _, err := utils.GenerateJWT(account.ToJWTAccInfo(), constant.TokenTypeRefresh)
	if err != nil {
		return nil, utils.ErrInternal("Failed generate jwt : "+err.Error(), "authUC.Login.GenerateJWT")
	}

	return &entity.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
