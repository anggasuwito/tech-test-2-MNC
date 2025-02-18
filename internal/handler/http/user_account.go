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

type UserAccountHandler struct {
	userAccUC usecase.UserAccUC
}

func NewUserAccHandler(
	userAccUC usecase.UserAccUC,
) *UserAccountHandler {
	return &UserAccountHandler{
		userAccUC: userAccUC,
	}
}

func (h *UserAccountHandler) SetupHandlers(r *gin.Engine) {
	userAccPathV1 := r.Group("/v1")
	userAccPathV1.Use(middleware.TokenChecker)
	userAccPathV1.POST("/user-account/update-profile", h.updateProfile)
}

func (h *UserAccountHandler) updateProfile(c *gin.Context) {
	var req entity.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "TransactionHandler.updateProfile.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ctxVal := middleware.GetContextValue(c)
	req.AccountID = ctxVal.AccountInfo.ID
	resp, err := h.userAccUC.UpdateProfile(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}
