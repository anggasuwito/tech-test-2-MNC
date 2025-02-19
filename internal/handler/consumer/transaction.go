package consumer

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/usecase"
	"tech-test-2-MNC/internal/utils"
	nsqd "tech-test-2-MNC/pkg/nsq"
	"time"
)

type TransactionHandler interface {
	UpdateTransactionStatus(msg *nsq.Message) error
}

type transactionHandler struct {
	nsq           *nsqd.NSQ
	transactionUC usecase.TransactionUC
}

func NewTransactionHandler(client *nsqd.NSQ, tuc usecase.TransactionUC) TransactionHandler {
	return &transactionHandler{
		nsq:           client,
		transactionUC: tuc,
	}
}

func (h *transactionHandler) UpdateTransactionStatus(msg *nsq.Message) error {
	var param entity.FinishTransactionMessage
	requeueTime := h.nsq.GetConsumerRequeueTime(constant.ConsumerUpdateTransactionStatus)
	maxAttempts := h.nsq.GetConsumerMaxAttempts(constant.ConsumerUpdateTransactionStatus)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		err = utils.ErrInternal("failed unmarshal : "+err.Error(), "transactionHandler.UpdateTransactionStatus.Unmarshal")
		utils.ResponseErrorNSQ(constant.ConsumerUpdateTransactionStatus, err)
		msg.RequeueWithoutBackoff(requeueTime)
		return err
	}

	if maxAttempts != 0 && msg.Attempts > maxAttempts {
		err = utils.ErrInternal("max attempt reached", "transactionHandler.UpdateTransactionStatus.maxAttempts")
		utils.ResponseErrorNSQ(constant.ConsumerUpdateTransactionStatus, err)
		msg.Finish()
		return nil
	}

	_, err = h.transactionUC.UpdateTransactionStatus(ctx, &entity.UpdateTransactionStatusRequest{
		TransactionID: param.TransactionID,
	})
	if err != nil {
		utils.ResponseErrorNSQ(constant.ConsumerUpdateTransactionStatus, err)
		msg.RequeueWithoutBackoff(requeueTime)
		return err
	}

	msg.Finish()
	return nil
}
