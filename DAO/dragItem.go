package dao

import (
	"context"
	"time"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

// 添加带事务的删除方法
func DeleteDraggableItemTx(ctx context.Context, tx *gorm.DB, itemID int64) error {
	return tx.WithContext(ctx).
		Model(&model.DraggableItem{}).
		Where("id = ?", itemID).
		Update("delete_time", time.Now().Unix()).Error
}

func UpdateDraggableItemTx(ctx context.Context, tx *gorm.DB, item *model.DraggableItem) error {
	return tx.WithContext(ctx).
		Model(&model.DraggableItem{}).
		Select("content", "week_type", "classroom", "update_time").
		Where("id = ? AND delete_time = 0", item.ID).
		Updates(item).
		Error
}

func GetClassNamesByItemID(ctx context.Context, itemID int64) ([]string, error) {
	var classNames []string
	err := mysql.GetDB().WithContext(ctx).
		Table("draggable_class_sheet dcs").
		Joins("JOIN class c ON dcs.class_id = c.id").
		Where("dcs.item_id = ?", itemID).
		Pluck("c.name", &classNames).Error

	return classNames, err
}

func DeleteItemClassRelationsTx(ctx context.Context, tx *gorm.DB, itemID int64) error {
	return tx.WithContext(ctx).
		Where("item_id = ?", itemID).
		Delete(&model.DraggableClassSheet{}).Error
}

// 创建元素-班级关联
func CreateItemSheetRelationTx(ctx context.Context, tx *gorm.DB, itemID, classID int64) error {
	relation := &model.DraggableClassSheet{
		ItemID:  itemID,
		ClassID: classID,
	}
	return tx.WithContext(ctx).Create(relation).Error
}

func CreateDraggableItemTx(ctx context.Context, tx *gorm.DB, item *model.DraggableItem) error {
	return tx.WithContext(ctx).Create(item).Error
}

func ListDraggableItemsByClass(ctx context.Context, classID int64) ([]*model.DraggableItem, error) {
	var items []*model.DraggableItem
	err := mysql.GetDB().WithContext(ctx).
		Joins("JOIN draggable_class_sheet ON draggable_class_sheet.item_id = draggable_item.id").
		Where("draggable_class_sheet.class_id = ? AND draggable_item.delete_time = 0", classID).
		Find(&items).Error
	return items, err
}

func CreateDraggableItem(ctx context.Context, item *model.DraggableItem) error {
	return mysql.GetDB().WithContext(ctx).Create(item).Error
}

func ListDraggableItemsBySheet(ctx context.Context, sheetID int64) ([]*model.DraggableItem, error) {
	var items []*model.DraggableItem
	err := mysql.GetDB().WithContext(ctx).
		Where("sheet_id = ? AND delete_time = 0", sheetID).
		Find(&items).Error
	return items, err
}

func UpdateDraggableItem(ctx context.Context, item *model.DraggableItem) error {
	// 显式指定更新条件和更新字段
	return mysql.GetDB().WithContext(ctx).
		Model(&model.DraggableItem{}).
		Where("id = ?", item.ID). // 同时校验所属sheet
		Updates(map[string]interface{}{
			"content":     item.Content,
			"update_time": time.Now(),
		}).Error
}

func GetDraggableItemByID(ctx context.Context, id int64) (*model.DraggableItem, error) {
	var item model.DraggableItem
	err := mysql.GetDB().WithContext(ctx).
		Where("id = ? AND delete_time = 0", id).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return &item, nil
}

func DeleteDraggableItem(ctx context.Context, itemID, sheetID int64) error {
	return mysql.GetDB().WithContext(ctx).
		Model(&model.DraggableItem{}).
		Where("id = ? AND sheet_id = ?", itemID, sheetID).
		Update("delete_time", time.Now().Unix()).Error
}

// GetCellByDraggableItemIDTx 根据拖拽元素ID获取所在单元格记录
func GetCellByDraggableItemIDTx(ctx context.Context, tx *gorm.DB, sheetID, DraggableItemID int64) (*model.Cell, error) {
	var cell model.Cell
	err := tx.WithContext(ctx).
		Where("sheet_id = ? AND item_id = ? AND delete_time = 0", sheetID, DraggableItemID).
		First(&cell).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cell, nil
}

// GetCellByPositionTx 根据工作表ID及行、列坐标获取单元格记录
func GetCellByPositionTx(ctx context.Context, tx *gorm.DB, sheetID int64, rowIndex, colIndex int) (*model.Cell, error) {
	var cell model.Cell
	err := tx.WithContext(ctx).
		Where("sheet_id = ? AND row_index = ? AND col_index = ? AND delete_time = 0", sheetID, rowIndex, colIndex).
		First(&cell).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cell, nil
}

// GetDraggableItemByIDTx 根据工作表 ID 和拖拽元素 ID 获取待拖拽元素记录
func GetDraggableItemByIDTx(ctx context.Context, tx *gorm.DB, DraggableItemID int64) (*model.DraggableItem, error) {
	var item model.DraggableItem
	err := tx.WithContext(ctx).
		Where("id = ? AND delete_time = 0", DraggableItemID).
		First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
