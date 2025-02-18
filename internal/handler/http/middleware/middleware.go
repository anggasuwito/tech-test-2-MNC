package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"tech-test-2-MNC/internal/domain/entity"
	"tech-test-2-MNC/internal/utils"
)

const (
	dataKey = "data"
)

func TokenChecker(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	token := strings.Split(authorization, "Bearer ")
	if len(token) < 2 || token[1] == "" {
		utils.ResponseError(c, utils.ErrUnauthorized("You are not authorized", "middleware.TokenChecker"))
		return
	}
	tokenStr := token[1]
	tokenClaims, err := utils.VerifyJWT(tokenStr)
	if err != nil {
		utils.ResponseError(c, utils.ErrUnauthorized("Invalid Token", "middleware.TokenChecker.VerifyJWT"))
		return
	}

	c.Set(dataKey, tokenClaims)
	c.Next()
}

func GetContextValue(c *gin.Context) entity.JWTClaim {
	if value, exists := c.Get(dataKey); exists {
		if data, ok := value.(entity.JWTClaim); ok {
			return data
		}
	}
	return entity.JWTClaim{}
}
