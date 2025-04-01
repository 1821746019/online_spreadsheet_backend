package model

import (
	"time"

	"gorm.io/gorm"
)

// Cell 单元格数据表
type Cell struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	SheetID    int64          `gorm:"column:sheet_id;not null;comment:所属工作表ID" json:"sheet_id"`
	RowIndex   int            `gorm:"column:row_index;not null;comment:行号（从1开始）" json:"row_index"`
	ColIndex   int            `gorm:"column:col_index;not null;comment:列号（从1开始）" json:"col_index"`
	Content    string         `gorm:"column:content;type:text;comment:单元格内容" json:"content"`
	ItemID     *int64         `gorm:"column:item_id;comment:关联的可拖放元素ID" json:"item_id"` // 使用指针表示可空
	LastEditBy int64          `gorm:"column:last_edit_by;not null;comment:最后编辑者ID" json:"last_edit_by"`
	Version    int64          `gorm:"column:version;default:0;comment:乐观锁版本号" json:"version"` // 乐观锁版本号，用于处理并发冲突
	CreateTime time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;default:0;comment:逻辑删除时间戳" json:"delete_time"`
}

// DragItem 可拖放元素表
type DragItem struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	SheetID    int64          `gorm:"column:sheet_id;not null;comment:工作表ID" json:"sheet_id"`
	Content    string         `gorm:"column:content;type:text;not null;comment:元素内容" json:"content"`
	CreatorID  int64          `gorm:"column:creator_id;not null;comment:创建者ID" json:"creator_id"`
	CreateTime time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;default:0;comment:逻辑删除时间戳" json:"delete_time"`
}
