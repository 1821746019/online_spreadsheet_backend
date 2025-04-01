package dao

import (
	"context"
	"fmt"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
	"gorm.io/gorm"
)

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
func UpdateCell(ctx context.Context, sheetID int64, content string, rowIndex, colIndex int, userID int64) error {
	return mysql.GetDB().WithContext(ctx).Model(&model.Cell{}).
		Where("sheet_id = ? AND row_index = ? AND col_index = ? AND delete_time = 0", sheetID, rowIndex, colIndex).
		Updates(map[string]interface{}{
			"content":     content,
			"update_time": gorm.Expr("NOW()"),
			"version":     gorm.Expr("version + 1"),
			"update_by":   userID,
		}).Error
}

func UpdateCellWithVersion(ctx context.Context, sheetID int64, value string, row, column int, version int64, userID int64) error {
	result := mysql.GetDB().WithContext(ctx).Exec(`
    UPDATE cells 
    SET content = ?, version = ?, update_by = ?, update_time = NOW()
    WHERE sheet_id = ? AND row_index = ? AND col_index = ? AND version = ?-1`,
		value, version, userID, sheetID, row, column, version)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("optimistic lock failed")
	}
	return nil
}

func GetCellWithVersion(ctx context.Context, sheetID int64, row, column int) (*model.Cell, error) {
	var cell model.Cell
	err := mysql.GetDB().WithContext(ctx).
		Where("sheet_id = ? AND row = ? AND column = ?", sheetID, row, column).
		First(&cell).Error
	return &cell, err
}
