package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"tech-test-2-MNC/config"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/domain/entity"
	"time"
)

const claimsDataKey = "claims_data"

func GenerateJWT(acc *entity.JWTClaimAccountInfo, tokenType string) (tokenStr string, data entity.JWTClaim, err error) {
	var (
		tokenB          = jwt.New(jwt.SigningMethodHS256)
		claims          = tokenB.Claims.(jwt.MapClaims)
		cfg             = config.GetConfig()
		expiredDuration string
		secretKey       string
	)

	switch tokenType {
	case constant.TokenTypeAccess:
		expiredDuration = cfg.AccessTokenExpireDuration
		secretKey = cfg.AccessTokenSecret
	case constant.TokenTypeRefresh:
		expiredDuration = cfg.RefreshTokenExpireDuration
		secretKey = cfg.RefreshTokenSecret
	}

	tokenExpiredDuration, _ := time.ParseDuration(expiredDuration)

	// Set payload
	expiredAt := TimeNow().Add(tokenExpiredDuration).Unix()
	data = entity.JWTClaim{
		ID:          uuid.New().String(),
		ExpiredAt:   expiredAt,
		AccountInfo: acc,
	}

	claims["expired_at"] = expiredAt
	claims[claimsDataKey] = data

	tokenStr, err = tokenB.SignedString([]byte(secretKey))
	if err != nil {
		return "", data, err
	}
	return tokenStr, data, err
}

func VerifyJWT(token, tokenType string) (jwtClaim entity.JWTClaim, err error) {
	var (
		cfg    = config.GetConfig()
		secret string
	)

	switch tokenType {
	case constant.TokenTypeAccess:
		secret = cfg.AccessTokenSecret
	case constant.TokenTypeRefresh:
		secret = cfg.RefreshTokenSecret
	}

	tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return jwtClaim, err
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return jwtClaim, err
	}

	jsonClaims, err := json.Marshal(claims[claimsDataKey])
	if err != nil {
		return jwtClaim, err
	}

	err = json.Unmarshal(jsonClaims, &jwtClaim)
	if err != nil {
		return jwtClaim, err
	}
	return jwtClaim, nil
}
