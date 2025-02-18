package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/usecase"
	"tech-test-2-MNC/internal/utils"
	"time"
)

type AuthHandler struct {
	authUC usecase.AuthUC
}

func NewAuthHandler(
	authUC usecase.AuthUC,
) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) SetupHandlers(r *gin.Engine) {
	authPathV1 := r.Group("/v1")
	authPathV1.POST("/auth/register", h.register)
	authPathV1.POST("/auth/login", h.login)
}

func (h *AuthHandler) register(c *gin.Context) {
	var req entity.AuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "AuthHandler.register.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	resp, err := h.authUC.Register(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}

func (h *AuthHandler) login(c *gin.Context) {
	var req entity.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.ErrBadRequest("Invalid Body : "+err.Error(), "AuthHandler.login.ShouldBindJSON"))
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	resp, err := h.authUC.Login(ctx, &req)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, "", resp)
}
