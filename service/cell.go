package service

import (
	"context"

	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"go.uber.org/zap"
)

func GetCells(ctx context.Context, userID, sheetID int64) ([]DTO.CellDTO, *apiError.ApiError) {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("GetSheet 查询权限失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限读取该工作表"}
	}
	// 查询单元格
	cells, err := dao.GetCellsBySheetID(ctx, sheetID)
	if err != nil {
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "获取单元格失败"}
	}
	// 转换为 DTO
	var result []DTO.CellDTO
	for _, c := range cells {
		result = append(result, DTO.CellDTO{
			ID:       c.ID,
			SheetID:  c.SheetID,
			RowIndex: int(c.RowIndex),
			ColIndex: int(c.ColIndex),
			ItemID:   c.ItemID,
		})
	}
	return result, nil
}

// // 更新单元格
// func UpdateCell(ctx context.Context, userID, sheetID int64, req DTO.UpdateCellRequestDTO) *apiError.ApiError {
// 	perm, err := dao.GetPermission(ctx, userID, sheetID)
// 	if err != nil {
// 		zap.L().Error("GetSheet 查询权限失败", zap.Error(err))
// 		return &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
// 	}
// 	if perm == nil {
// 		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
// 	}
// 	if err := dao.UpdateCell(ctx, sheetID, req.Content, req.RowIndex, req.ColIndex, userID); err != nil {
// 		zap.L().Error("UpdateSheet 更新失败", zap.Error(err))
// 		return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
// 	}
// 	return nil
// }
