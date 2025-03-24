package dao

import (
	"context"
	"time"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

func CreateDraggableItem(ctx context.Context, item *model.DragItem) error {
	return mysql.GetDB().WithContext(ctx).Create(item).Error
}

func ListDraggableItemsBySheet(ctx context.Context, sheetID int64) ([]*model.DragItem, error) {
	var items []*model.DragItem
	err := mysql.GetDB().WithContext(ctx).
		Where("sheet_id = ? AND delete_time = 0", sheetID).
		Find(&items).Error
	return items, err
}

func UpdateDraggableItem(ctx context.Context, item *model.DragItem) error {
	// 显式指定更新条件和更新字段
	return mysql.GetDB().WithContext(ctx).
		Model(&model.DragItem{}).
		Where("id = ? AND sheet_id = ?", item.ID, item.SheetID). // 同时校验所属sheet
		Updates(map[string]interface{}{
			"content":     item.Content,
			"update_time": time.Now(),
		}).Error
}

func GetDraggableItemByID(ctx context.Context, id int64) (*model.DragItem, error) {
	var item model.DragItem
	err := mysql.GetDB().WithContext(ctx).
		Where("id = ? AND delete_time = 0", id).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func DeleteDraggableItem(ctx context.Context, itemID, sheetID int64) error {
	return mysql.GetDB().WithContext(ctx).
		Model(&model.DragItem{}).
		Where("id = ? AND sheet_id = ?", itemID, sheetID).
		Update("delete_time", time.Now().Unix()).Error
}

// GetCellByDragItemIDTx 根据拖拽元素ID获取所在单元格记录
func GetCellByDragItemIDTx(ctx context.Context, tx *gorm.DB, sheetID, dragItemID int64) (*model.Cell, error) {
	var cell model.Cell
	err := tx.WithContext(ctx).
		Where("sheet_id = ? AND item_id = ? AND delete_time = 0", sheetID, dragItemID).
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

// UpdateCellTx 更新单元格记录
func UpdateCellTx(ctx context.Context, tx *gorm.DB, cell *model.Cell) error {
	return tx.WithContext(ctx).Save(cell).Error
}

// GetDraggableItemByIDTx 根据工作表 ID 和拖拽元素 ID 获取待拖拽元素记录
func GetDraggableItemByIDTx(ctx context.Context, tx *gorm.DB, sheetID, dragItemID int64) (*model.DragItem, error) {
	var item model.DragItem
	err := tx.WithContext(ctx).
		Where("sheet_id = ? AND id = ? AND delete_time = 0", sheetID, dragItemID).
		First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
