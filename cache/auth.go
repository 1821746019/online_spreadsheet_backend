package cache

import (
	"context"
	"time"

	"github.com/sztu/mutli-table/DAO/Redis"
)

// AddTokenToBlacklist 将token添加到Redis黑名单中，并指定过期时间。
// 它使用提供的token和预定义模板生成Redis键，然后将此键设置为值"1"并指定过期时间。
//
// 参数:
//   - ctx: 操作的上下文，允许取消和超时控制。
//   - token: 要加入黑名单的token。
//   - expiration: token在黑名单中应保留的时间。
//
// 返回:
//   - error: 如果操作失败，则返回错误对象，否则返回nil。
func AddTokenToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	err := Redis.GetRedisClient().Set(ctx, key, "1", expiration).Err()

	return err
}
