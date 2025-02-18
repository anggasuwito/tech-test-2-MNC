package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/handler/http/middleware"
	"tech-test-2-MNC/internal/usecase"
	"tech-test-2-MNC/internal/utils"
	"time"
)

type TransactionHandler struct {
	transactionUC usecase.TransactionUC
}

func NewTransactionHandler(
	transactionUC usecase.TransactionUC,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUC: transactionUC,
	}
}

func (h *TransactionHandler) SetupHandlers(r *gin.Engine) {
	transactionPathV1 := r.Group("/v1")
	transactionPathV1.Use(middleware.TokenChecker)
	transactionPathV1.POST("/transaction/topup", h.topup)
	transactionPathV1.POST("/transaction/transfer", h.transfer)
	transactionPathV1.POST("/transaction/payment", h.payment)
	transactionPathV1.GET("/transaction/report", h.report)
}

func (h *TransactionHandler) topup(c *gin.Context) {
	var req entity.TransactionTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "TransactionHandler.topup.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ctxVal := middleware.GetContextValue(c)
	req.AccountID = ctxVal.AccountInfo.ID
	resp, err := h.transactionUC.Topup(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}

func (h *TransactionHandler) transfer(c *gin.Context) {
	var req entity.TransactionTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "TransactionHandler.transfer.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ctxVal := middleware.GetContextValue(c)
	req.AccountID = ctxVal.AccountInfo.ID
	resp, err := h.transactionUC.Transfer(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}

func (h *TransactionHandler) payment(c *gin.Context) {
	var req entity.TransactionPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "TransactionHandler.payment.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ctxVal := middleware.GetContextValue(c)
	req.AccountID = ctxVal.AccountInfo.ID
	resp, err := h.transactionUC.Payment(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}

func (h *TransactionHandler) report(c *gin.Context) {
	var req entity.TransactionReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "TransactionHandler.report.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ctxVal := middleware.GetContextValue(c)
	req.Search = append(req.Search, &entity.Filter{
		Field: "account_id",
		Value: ctxVal.AccountInfo.ID,
	})
	resp, err := h.transactionUC.Report(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}
