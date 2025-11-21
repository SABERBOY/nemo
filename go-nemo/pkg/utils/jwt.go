package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	ExpireTime              = 15 * 60 * 1000 // 15 minutes in milliseconds
	DefaultMinSSTokenExpire = 7200000        // 2 hours in milliseconds
)

type UserClaim struct {
	UserUuid string `json:"userUuid"`
	AppKey   string `json:"appKey"`
	AppId    string `json:"appId"`
	jwt.RegisteredClaims
}

func CreateToken(appKey, userUuid, appId, secret string) (string, error) {
	claims := UserClaim{
		UserUuid: userUuid,
		AppKey:   appKey,
		AppId:    appId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ExpireTime) * time.Millisecond)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Set headers to match Java: typ=JWT, alg=HS256 (default)
	return token.SignedString([]byte(secret))
}

func CreateSSToken(appKey, userUuid, appId, secret string, expireDuration int64) (string, int64, error) {
	expireTime := time.Now().Add(time.Duration(expireDuration) * time.Millisecond)
	claims := UserClaim{
		UserUuid: userUuid,
		AppKey:   appKey,
		AppId:    appId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, expireTime.UnixMilli(), err
}

func VerifyToken(secret, tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
