package model

import (
	"time"

	"gorm.io/gorm"
)

type Sheet struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增主键" json:"id"`                    // 自增主键
	Name       string         `gorm:"column:name;type:varchar(255);not null;comment:工作表名称" json:"name"`                  // 工作表名称
	CreatorID  int64          `gorm:"column:creator_id;not null;index;comment:创建者ID（关联user.id）" json:"creator_id"`       // 创建者ID（关联 user.id）
	Row        int            `gorm:"column:row;not null;comment:行数" json:"row"`                                         // 行数
	Col        int            `gorm:"column:col;not null;comment:列数" json:"col"`                                         // 列数
	CreateTime time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:记录的创建时间" json:"create_time"`   // 记录的创建时间
	UpdateTime time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:记录的最后更新时间" json:"update_time"` // 记录的最后更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;default:0;comment:逻辑删除时间戳" json:"delete_time"`                   // 逻辑删除时间戳（0 表示未删除）
}
