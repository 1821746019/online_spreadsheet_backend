package dao

import (
	"context"
	"errors"
	"time"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

// 根据班级ID获取班级信息
func GetClassByID(ctx context.Context, classID int64) (*model.Class, error) {
	var class model.Class
	err := mysql.GetDB().WithContext(ctx).
		Where("id = ? AND delete_time = 0", classID).
		First(&class).Error
	return &class, err
}

func CreateClass(ctx context.Context, class *model.Class) error {
	return mysql.GetDB().WithContext(ctx).Create(class).Error
}

func ClassNameExists(ctx context.Context, name string) (bool, error) {
	var exist bool
	err := mysql.GetDB().WithContext(ctx).
		Model(&model.Class{}).
		Select("1").
		Where("name = ? AND delete_time = 0", name).
		Limit(1).
		Scan(&exist).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return exist, nil
}

func ListClasses(ctx context.Context, userID int64, page, pageSize int) ([]*model.Class, int64, error) {
	var classes []*model.Class
	var total int64

	db := mysql.GetDB().WithContext(ctx).Model(&model.Class{}).
		Where("user_id =? AND delete_time = 0", userID) // 移除权限关联，直接使用班级ID

	// 查询总记录数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&classes).Error; err != nil {
		return nil, total, err
	}
	return classes, total, nil
}

func UpdateClassTx(ctx context.Context, tx *gorm.DB, class *model.Class) error {
	return tx.WithContext(ctx).Model(class).Updates(class).Error
}

func DeleteClassTx(ctx context.Context, tx *gorm.DB, classID int64) error {
	return tx.WithContext(ctx).Model(&model.Class{}).Where("id =? AND delete_time = 0", classID).Update("delete_time", time.Now().Unix()).Error
}

func GetClassByName(ctx context.Context, name string) (*model.Class, error) {
	var class model.Class
	err := mysql.GetDB().WithContext(ctx).
		Where("name =? AND delete_time = 0", name).
		First(&class).Error
	return &class, err
}

func UpdateClass(ctx context.Context, class *model.Class) error {
	return mysql.GetDB().WithContext(ctx).Model(class).Updates(class).Error
}

func DeleteClass(ctx context.Context, classID int64) error {
	return mysql.GetDB().WithContext(ctx).Model(&model.Class{}).Where("id =? AND delete_time = 0", classID).Update("delete_time", time.Now().Unix()).Error
}
