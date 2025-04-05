package dao

import (
	"context"
	"time"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

// CreateSheetTx 使用事务插入一条 Sheet 记录
func CreateSheetTx(ctx context.Context, tx *gorm.DB, sheet *model.Sheet) error {
	return tx.WithContext(ctx).Model(model.Sheet{}).Create(sheet).Error
}

// ListSheets
func ListSheets(ctx context.Context, userID, classID int64, page, pageSize int) ([]*model.Sheet, int64, error) {
	var sheets []*model.Sheet
	var total int64

	db := mysql.GetDB().WithContext(ctx).Model(&model.Sheet{}).
		Where("class_id = ? AND delete_time = ?", classID, 0) // 移除权限关联，直接使用班级ID

	// 查询总记录数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&sheets).Error; err != nil {
		return nil, total, err
	}
	return sheets, total, nil
}

// GetSheetByID 根据工作表 ID 查询记录
func GetSheetByID(ctx context.Context, sheetID int64) (*model.Sheet, error) {
	var sheet model.Sheet
	err := mysql.GetDB().WithContext(ctx).Where("id = ? AND delete_time = ?", sheetID, 0).First(&sheet).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sheet, err
}

// UpdateSheet 更新工作表记录
func UpdateSheet(ctx context.Context, sheet *model.Sheet) error {
	return mysql.GetDB().WithContext(ctx).Model(sheet).Updates(sheet).Error
}

// DeleteSheet 逻辑删除工作表，更新 delete_time 字段为当前时间戳
func DeleteSheet(ctx context.Context, sheetID int64) error {
	return mysql.GetDB().WithContext(ctx).
		Model(&model.Sheet{}).
		Where("id = ? AND delete_time = ?", sheetID, 0).
		Update("delete_time", time.Now().Unix()).Error
}

func GetSheetByClassIDandWeek(ctx context.Context, classID int64, week int) (*model.Sheet, error) {
	var sheet model.Sheet
	err := mysql.GetDB().WithContext(ctx).Where("class_id =? AND week =? AND delete_time =?", classID, week, 0).First(&sheet).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sheet, err
}
