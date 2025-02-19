package usecase

import (
	"context"
	"database/sql"
	"math"
	"tech-test-2-MNC/config"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/domain/model"
	"tech-test-2-MNC/internal/repository"
	"tech-test-2-MNC/internal/utils"
)

type TransactionUC interface {
	Topup(ctx context.Context, req *entity.TransactionTopupRequest) (*entity.TransactionTopupResponse, error)
	Transfer(ctx context.Context, req *entity.TransactionTransferRequest) (*entity.TransactionTransferResponse, error)
	Payment(ctx context.Context, req *entity.TransactionPaymentRequest) (*entity.TransactionPaymentResponse, error)
	Report(ctx context.Context, req *entity.TransactionReportRequest) (*entity.TransactionReportResponse, error)
	UpdateTransactionStatus(ctx context.Context, req *entity.UpdateTransactionStatusRequest) (*entity.UpdateTransactionStatusResponse, error)
}

type transactionUC struct {
	producer        *config.NSQProducer
	txWrapper       repository.TransactionWrapper
	accRepo         repository.UserAccountRepo
	transactionRepo repository.TransactionRepo
}

func NewTransactionUC(
	producer *config.NSQProducer,
	txWrapper repository.TransactionWrapper,
	accRepo repository.UserAccountRepo,
	transactionRepo repository.TransactionRepo,
) TransactionUC {
	return &transactionUC{
		producer:        producer,
		txWrapper:       txWrapper,
		accRepo:         accRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *transactionUC) Topup(ctx context.Context, req *entity.TransactionTopupRequest) (*entity.TransactionTopupResponse, error) {
	var (
		transaction       = &model.Transaction{}
		transactionDetail = &model.TransactionDetail{}
	)

	account, err := u.accRepo.GetUserAccountByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	if err = u.txWrapper.ExecuteTransaction(ctx,
		func(ctxTX context.Context) error {
			//create transaction
			transaction.CreateNewTransaction(constant.TransactionCategoryTopUp, "", req.Amount)
			err = u.transactionRepo.CreateTransaction(ctxTX, transaction)
			if err != nil {
				return err
			}

			//create transaction details
			transactionDetail.CreateNewTransactionDetail(transaction.ID, req.AccountID, constant.TransactionTypeCredit, req.Amount, account.AccountBalance.Balance)
			err = u.transactionRepo.CreateTransactionDetail(ctxTX, transactionDetail)
			if err != nil {
				return err
			}

			//publish to topic finish_transaction
			if err = u.producer.Publish(constant.TopicFinishTransaction, entity.FinishTransactionMessage{
				TransactionID: transaction.ID,
			}); err != nil {
				return utils.ErrInternal("Failed to publish message : "+err.Error(), "transactionUC.Topup.producer.Publish")
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return &entity.TransactionTopupResponse{
		transaction.ToEntity(transactionDetail.Type, transactionDetail.BalanceBefore, transactionDetail.BalanceAfter),
	}, nil
}

func (u *transactionUC) Transfer(ctx context.Context, req *entity.TransactionTransferRequest) (*entity.TransactionTransferResponse, error) {
	var (
		transaction                  = &model.Transaction{}
		transactionDetailSource      = &model.TransactionDetail{}
		transactionDetailDestination = &model.TransactionDetail{}
	)

	account, err := u.accRepo.GetUserAccountByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	destinationAccount, err := u.accRepo.GetUserAccountByID(ctx, req.TargetAccountID)
	if err != nil {
		return nil, err
	}

	if req.Amount > account.AccountBalance.Balance {
		return nil, utils.ErrBadRequest("Balance is not enough", "transactionUC.Transfer.Balance")
	}

	if req.AccountID == req.TargetAccountID {
		return nil, utils.ErrBadRequest("Invalid transfer target", "transactionUC.Transfer.TargetAccountID")
	}

	if err = u.txWrapper.ExecuteTransaction(ctx,
		func(ctxTX context.Context) error {
			//create transaction
			transaction.CreateNewTransaction(constant.TransactionCategoryTransfer, "", req.Amount)
			err = u.transactionRepo.CreateTransaction(ctxTX, transaction)
			if err != nil {
				return err
			}

			//create transaction details source
			transactionDetailSource.CreateNewTransactionDetail(transaction.ID, req.AccountID, constant.TransactionTypeDebit, req.Amount, account.AccountBalance.Balance)
			err = u.transactionRepo.CreateTransactionDetail(ctxTX, transactionDetailSource)
			if err != nil {
				return err
			}

			//create transaction details destination
			transactionDetailDestination.CreateNewTransactionDetail(transaction.ID, req.TargetAccountID, constant.TransactionTypeCredit, req.Amount, destinationAccount.AccountBalance.Balance)
			err = u.transactionRepo.CreateTransactionDetail(ctxTX, transactionDetailDestination)
			if err != nil {
				return err
			}

			//publish to topic finish_transaction
			if err = u.producer.Publish(constant.TopicFinishTransaction, entity.FinishTransactionMessage{
				TransactionID: transaction.ID,
			}); err != nil {
				return utils.ErrInternal("Failed to publish message : "+err.Error(), "transactionUC.Transfer.producer.Publish")
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return &entity.TransactionTransferResponse{
		transaction.ToEntity(transactionDetailSource.Type, transactionDetailSource.BalanceBefore, transactionDetailSource.BalanceAfter),
	}, nil
}

func (u *transactionUC) Payment(ctx context.Context, req *entity.TransactionPaymentRequest) (*entity.TransactionPaymentResponse, error) {
	var (
		transaction       = &model.Transaction{}
		transactionDetail = &model.TransactionDetail{}
	)

	account, err := u.accRepo.GetUserAccountByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	if req.Amount > account.AccountBalance.Balance {
		return nil, utils.ErrBadRequest("Balance is not enough", "transactionUC.Payment.Balance")
	}

	if err = u.txWrapper.ExecuteTransaction(ctx,
		func(ctxTX context.Context) error {
			//create transaction
			transaction.CreateNewTransaction(constant.TransactionCategoryPayment, "", req.Amount)
			err = u.transactionRepo.CreateTransaction(ctxTX, transaction)
			if err != nil {
				return err
			}

			//create transaction details
			transactionDetail.CreateNewTransactionDetail(transaction.ID, req.AccountID, constant.TransactionTypeDebit, req.Amount, account.AccountBalance.Balance)
			err = u.transactionRepo.CreateTransactionDetail(ctxTX, transactionDetail)
			if err != nil {
				return err
			}

			//publish to topic finish_transaction
			if err = u.producer.Publish(constant.TopicFinishTransaction, entity.FinishTransactionMessage{
				TransactionID: transaction.ID,
			}); err != nil {
				return utils.ErrInternal("Failed to publish message : "+err.Error(), "transactionUC.Payment.producer.Publish")
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return &entity.TransactionPaymentResponse{
		transaction.ToEntity(transactionDetail.Type, transactionDetail.BalanceBefore, transactionDetail.BalanceAfter),
	}, nil
}

func (u *transactionUC) Report(ctx context.Context, req *entity.TransactionReportRequest) (*entity.TransactionReportResponse, error) {
	var (
		limit = req.Limit
		page  = req.Page
	)

	if limit > 100 || limit == 0 {
		limit = 100
	}

	if page < 1 {
		page = 1
	}

	data, totalData, err := u.transactionRepo.GetTransactionDetails(ctx, req.Search, req.Sort, page, limit)
	if err != nil {
		return nil, utils.ErrInternal("Failed get transaction report list : "+err.Error(), "transactionUC.Report.transactionRepo.GetTransactionDetails")
	}

	var resp []*entity.Transaction
	for _, d := range data {
		resp = append(resp, d.ToTransactionEntity())
	}

	return &entity.TransactionReportResponse{
		ListPaginationResponse: &entity.ListPaginationResponse{
			CurrentPage: page,
			TotalPage:   int64(math.Ceil(float64(totalData) / float64(limit))),
			TotalData:   totalData,
			PerPage:     limit,
		},
		Data: resp,
	}, nil
}

func (u *transactionUC) UpdateTransactionStatus(ctx context.Context, req *entity.UpdateTransactionStatusRequest) (*entity.UpdateTransactionStatusResponse, error) {
	transaction, err := u.transactionRepo.GetTransactionByID(ctx, req.TransactionID)
	if err != nil {
		return nil, err
	}

	if transaction.Status != constant.TransactionStatusPending {
		return &entity.UpdateTransactionStatusResponse{}, nil
	}

	if err = u.txWrapper.ExecuteTransaction(ctx,
		func(ctxTX context.Context) error {
			//update transaction status
			transaction.Status = constant.TransactionStatusSuccess
			transaction.CompletedAt = sql.NullTime{Time: utils.TimeNow(), Valid: true}
			err = u.transactionRepo.UpdateTransaction(ctxTX, transaction)
			if err != nil {
				return err
			}

			for _, detail := range transaction.TransactionDetails {
				//get and lock account balance
				accountBalance, err := u.accRepo.GetAndLockAccountBalance(ctxTX, detail.AccountID)
				if err != nil {
					return err
				}

				//update account balance
				accountBalance.Balance = detail.BalanceAfter
				err = u.accRepo.UpdateAccountBalance(ctxTX, accountBalance)
				if err != nil {
					return err
				}
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}
