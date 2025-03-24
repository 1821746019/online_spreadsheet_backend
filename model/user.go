package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户信息表：存储用户基本信息及状态
type User struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增主键，唯一标识用户记录" json:"id"`           // 自增主键，唯一标识用户记录
	UserID     int64          `gorm:"column:user_id;not null;comment:用户ID，用于业务中的用户唯一标识" json:"user_id"`                  // 用户ID，用于业务中的用户唯一标识
	Username   string         `gorm:"column:username;not null;comment:用户名，唯一且不区分大小写" json:"username"`                    // 用户名，唯一且不区分大小写
	Password   string         `gorm:"column:password;not null;comment:用户密码，存储的是哈希值" json:"password"`                     // 用户密码，存储的是哈希值
	Email      string         `gorm:"column:email;comment:用户邮箱，可为空" json:"email"`                                        // 用户邮箱，可为空
	CreateTime time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:记录的创建时间" json:"create_time"`   // 记录的创建时间
	UpdateTime time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:记录的最后更新时间" json:"update_time"` // 记录的最后更新时间
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:逻辑删除时间，NULL表示未删除" json:"delete_time"`                    // 逻辑删除时间，NULL表示未删除
}
