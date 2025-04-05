// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePermission = "permission"

// Permission 权限控制表
type Permission struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`    // 用户ID
	SheetID    int64     `gorm:"column:sheet_id;not null;comment:工作表ID" json:"sheet_id"` // 工作表ID
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
	DeletedTime    gorm.DeletedAt     `gorm:"column:deleted_time" json:"deleted_time"`
}

// TableName Permission's table name
func (*Permission) TableName() string {
	return TableNamePermission
}
