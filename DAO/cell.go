package dao

import (
	"context"

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
func UpdateCell(ctx context.Context, sheetID int64, content string, rowIndex, colIndex int) error {
	return mysql.GetDB().WithContext(ctx).Model(&model.Cell{}).
		Where("sheet_id = ? AND row_index = ? AND col_index = ? AND delete_time = 0", sheetID, rowIndex, colIndex).
		Updates(map[string]interface{}{
			"content":     content,
			"update_time": gorm.Expr("NOW()"),
		}).Error
}
