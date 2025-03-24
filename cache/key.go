package cache

import "fmt"

const (
	BlackListTokenKeyTemplate = "blacklist:token:%v"
)

// GenerateRedisKey 通过格式化给定的模板字符串和提供的参数生成一个 Redis key。
//
// 参数:
//   - template: 一个包含格式占位符的字符串模板。
//   - param: 一个可变参数，表示要格式化到模板中的值。
//
// 返回值:
//   - string: 一个格式化的字符串，可以用作 Redis key。
func GenerateRedisKey(template string, param ...any) string {
	return fmt.Sprintf(template, param...)
}
