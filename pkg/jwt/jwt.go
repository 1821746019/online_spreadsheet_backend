package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySecret = "111"

const (
	// AccessTokenName 是访问令牌的key
	AccessTokenName = "access"
	// RefreshTokenName 是刷新令牌的key
	RefreshTokenName = "refresh"
)

type MyClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 accessToken
func GenerateToken[T int64 | string | uint](userID T, username string) (accessToken string, err error) {
	var int64UserID int64
	switch v := any(userID).(type) {
	case uint:
		int64UserID = int64(v)
	case int64:
		int64UserID = v
	case string:
		// 尝试将 string 转为 uint
		parsedID, err := strconv.ParseUint(v, 10, 32) // 假设 uint 是 32 位
		if err != nil {
			return "", fmt.Errorf("invalid userID format, could not convert to uint: %v", err)
		}
		int64UserID = int64(parsedID)
	default:
		return "", fmt.Errorf("unsupported userID type")
	}

	// 定义生成 token 的闭包函数
	generate := func(userID int64, username string, tokenType string, validTime time.Duration) (string, error) {
		claims := MyClaims{
			UserID:    userID,
			Username:  username,
			TokenType: tokenType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(validTime)),
				Issuer:    "Ethen",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(mySecret))
	}

	accessToken, err = generate(int64UserID, username, AccessTokenName, time.Hour*24*7)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
