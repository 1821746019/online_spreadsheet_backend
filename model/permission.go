package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限控制表：管理用户对工作表的访问权限
type Permission struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增主键" json:"id"`
	UserID      int64          `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`
	SheetID     int64          `gorm:"column:sheet_id;not null;comment:工作表ID" json:"sheet_id"`
	AccessLevel string         `gorm:"column:access_level;type:enum('READ','WRITE','ADMIN');not null;default:'READ';comment:访问级别" json:"access_level"`
	CreateTime  time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeleteTime  gorm.DeletedAt `gorm:"column:delete_time;default:0;comment:逻辑删除时间戳" json:"delete_time"`
}
