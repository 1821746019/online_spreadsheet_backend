package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sztu/mutli-table/DAO/Redis"
	"github.com/sztu/mutli-table/cache"
	"github.com/sztu/mutli-table/pkg/jwt"
	"go.uber.org/zap"
)

const (
	// ContextUserIDKey 是上下文中用户ID的key
	ContextUserIDKey = "user_id"
	// ContextUsernameKey 是上下文中用户名的key
	ContextUsernameKey = "username"
)

// JWTAuthMiddleware 是一个 Gin 的中间件函数, 用于处理 JWT 认证。
// 它检查请求的 Authorization 头中是否存在有效的 JWT token。
// 如果 token 缺失、格式错误或无效, 它会返回未授权错误并中止请求。
// 如果 token 有效, 它会从 token 中提取用户信息并设置到请求上下文中。
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ResponseUnAuthorized(c, "请求未携带 token")
			zap.L().Info("请求未携带 token")
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ResponseUnAuthorized(c, "token 格式错误, 应为 `Bearer {token}`")
			zap.L().Info("token 格式错误")
			c.Abort()
			return
		}
		token := parts[1]
		myClaims, err := jwt.ParseToken(token)
		if err != nil {
			ResponseUnAuthorized(c, "token 解析失败")
			zap.L().Error("token 解析失败", zap.Error(err))
			c.Abort()
			return
		}
		if myClaims.TokenType != jwt.AccessTokenName {
			ResponseUnAuthorized(c, "token 类型错误")
			zap.L().Info("token 类型错误")
			c.Abort()
			return
		}
		if err := CheckTokenBlacklist(context.Background(), token); err != nil {
			ResponseUnAuthorized(c, "token 已被拉黑")
			zap.L().Info("token 已被拉黑")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, myClaims.UserID)
		c.Set(ContextUsernameKey, myClaims.Username)
		c.Next()
		return
	}
}

func CheckTokenBlacklist(ctx context.Context, token string) error {
	key := cache.GenerateRedisKey(cache.BlackListTokenKeyTemplate, token)
	err := Redis.GetRedisClient().Get(ctx, key).Err()
	if err == redis.Nil {
		// 没有查到，说明不在黑名单
		return nil
	}
	if err != nil {
		// 查询出错，返回原始错误
		return err
	}
	// 查到了，说明在黑名单
	return errors.New("token is blacklisted")
}
