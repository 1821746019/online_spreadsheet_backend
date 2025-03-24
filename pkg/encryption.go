package pkg

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/sztu/mutli-table/settings"
)

// EncryptPassword 用于加密密码
func EncryptPassword(password string) string {
	secret := settings.GetConfig().PasswordSecret
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
