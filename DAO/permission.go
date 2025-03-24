package dao

import (
	"context"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

// CreatePermissionTx 使用事务插入一条 Permission 记录
func CreatePermissionTx(tx *gorm.DB, permission *model.Permission) error {
	return tx.Create(permission).Error
}

// GetPermission 根据用户ID和工作表ID查询权限记录
func GetPermission(ctx context.Context, userID, sheetID int64) (*model.Permission, error) {
	var perm model.Permission
	err := mysql.GetDB().WithContext(ctx).
		Where("user_id = ? AND sheet_id = ? AND delete_time = ?", userID, sheetID, 0).
		First(&perm).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &perm, err
}
