package dao

import (
	"context"
	"errors"
	"time"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

func GetCellByDragItemIDTx(ctx context.Context, tx *gorm.DB, sheetID, dragItemID int64) (*model.Cell, error) {
	var cell model.Cell
	err := tx.WithContext(ctx).
		Where("sheet_id = ? AND item_id = ? AND delete_time = 0", sheetID, dragItemID).
		First(&cell).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 未找到视为正常情况
	}
	return &cell, err
}

func CountCellReferences(ctx context.Context, itemID int64) (int64, error) {
	var refCount int64
	err := mysql.GetDB().WithContext(ctx).
		Model(&model.Cell{}).
		Where("item_id = ? AND delete_time = 0", itemID).
		Count(&refCount).Error
	return refCount, err
}

func CreateBatchCellsTx(tx *gorm.DB, ctx context.Context, cells []model.Cell) error {
	if len(cells) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&cells).Error
}

// 获取单元格
func GetCellsBySheetID(ctx context.Context, sheetID int64) ([]model.Cell, error) {
	var cells []model.Cell
	err := mysql.GetDB().WithContext(ctx).Where("sheet_id = ? AND delete_time = 0", sheetID).Find(&cells).Error
	return cells, err
}

// 更新单元格
func UpdateCell(ctx context.Context, sheetID int64, cell *model.Cell) error {
	return mysql.GetDB().WithContext(ctx).Model(cell).
		Select("item_id", "last_modified_by", "version", "update_time").
		Updates(map[string]interface{}{
			"item_id":          cell.ItemID,
			"last_modified_by": cell.LastModifiedBy,
			"version":          cell.Version + 1,
			"update_time":      time.Now(),
		}).Error
}

// UpdateCellTx 更新单元格记录
func UpdateCellTx(ctx context.Context, tx *gorm.DB, cell *model.Cell) error {
	return tx.WithContext(ctx).Model(cell).
		Select("item_id", "last_modified_by", "version", "update_time").
		Updates(map[string]interface{}{
			"item_id":          cell.ItemID,
			"last_modified_by": cell.LastModifiedBy,
			"version":          cell.Version + 1,
			"update_time":      time.Now(),
		}).Error
}

func GetCellWithVersion(ctx context.Context, sheetID int64, row, column int) (*model.Cell, error) {
	var cell model.Cell
	err := mysql.GetDB().WithContext(ctx).
		Where("sheet_id = ? AND row = ? AND column = ?", sheetID, row, column).
		First(&cell).Error
	return &cell, err
}

func GetCellByPosition(ctx context.Context, sheetID int64, row, column int) (*model.Cell, error) {
	var cell model.Cell
	err := mysql.GetDB().WithContext(ctx).
		Where("sheet_id =? AND row_index =? AND col_index =? AND delete_time = 0", sheetID, row, column).
		First(&cell).Error
	return &cell, err
}
