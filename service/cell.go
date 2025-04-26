package service

import (
	"context"
	"time"

	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"go.uber.org/zap"
)

func GetCells(ctx context.Context, userID, sheetID int64) ([]DTO.CellDTO, *apiError.ApiError) {
	// 查询单元格
	cells, err := dao.GetCellsBySheetID(ctx, sheetID)
	if err != nil {
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "获取单元格失败"}
	}
	// 转换为 DTO
	var result []DTO.CellDTO
	for _, c := range cells {
		if c.ItemID == nil {
			result = append(result, DTO.CellDTO{
				ID:       c.ID,
				SheetID:  c.SheetID,
				RowIndex: int(c.RowIndex),
				ColIndex: int(c.ColIndex),
			})
		} else {
			dragItem, err := dao.GetDraggableItemByID(ctx, *c.ItemID)
			if err != nil {
				return nil, &apiError.ApiError{Code: code.ServerError, Msg: "获取拖拽项失败"}
			}
			result = append(result, DTO.CellDTO{
				ID:        c.ID,
				SheetID:   c.SheetID,
				RowIndex:  int(c.RowIndex),
				ColIndex:  int(c.ColIndex),
				ItemID:    c.ItemID,
				Content:   dragItem.Content,
				WeekType:  dragItem.WeekType,
				ClassRoom: dragItem.Classroom,
			})
		}
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

func DeleteItemInCell(ctx context.Context, userID, classID, sheetID int64, req DTO.DeleteItemInCellRequest) *apiError.ApiError {
	targetCell, err := dao.GetCellByPosition(ctx, sheetID, req.Row, req.Col)
	needToDelete := 0
	if err != nil {
		zap.L().Error("GetCellByRowAndCol 查询单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取单元格失败"}
	}
	if targetCell.ItemID == nil {
		return &apiError.ApiError{Code: code.InvalidParam, Msg: "删除的单元格为空"}
	}
	if targetCell.LastModifiedBy != userID {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该单元格"}
	}
	needToDelete = int(*targetCell.ItemID)
	targetCell.ItemID = nil
	targetCell.UpdateTime = time.Now()
	targetCell.LastModifiedBy = userID
	currentSheet, err := dao.GetSheetByID(ctx, sheetID)
	if err != nil || currentSheet == nil {
		zap.L().Error("获取工作表信息失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取工作表失败"}
	}

	// 更新当前单元格
	if err := dao.UpdateCell(ctx, sheetID, targetCell); err != nil {
		zap.L().Error("更新当前单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
	}

	// 获取该班级的所有工作表
	sheets, _, err := dao.ListSheets(ctx, userID, currentSheet.ClassID, 1, 1000) // 假设一个班级不会有超过1000个工作表
	if err != nil {
		zap.L().Error("获取班级工作表列表失败",
			zap.Int64("classID", currentSheet.ClassID),
			zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取班级工作表列表失败"}
	}

	// 遍历所有工作表，更新相同位置的单元格
	for _, sheet := range sheets {
		if sheet.ID == sheetID { // 跳过当前工作表（已处理）
			continue
		}

		// 获取目标单元格
		targetCell, err := dao.GetCellByPosition(ctx, sheet.ID, req.Row, req.Col)
		if err != nil || targetCell == nil {
			zap.L().Error("获取目标单元格失败",
				zap.Int64("sheetID", sheet.ID),
				zap.Error(err))
			continue
		}
		if targetCell.ItemID == nil {
			continue
		}
		if targetCell.LastModifiedBy != userID {
			continue
		}
		if int(*targetCell.ItemID) != needToDelete {
			continue
		}
		targetCell.ItemID = nil
		targetCell.UpdateTime = time.Now()
		targetCell.LastModifiedBy = userID

		// 更新目标单元格
		if err := dao.UpdateCell(ctx, sheet.ID, targetCell); err != nil {
			zap.L().Error("更新工作表单元格失败",
				zap.Int64("sheetID", sheet.ID),
				zap.Error(err))
		}
	}

	return nil
}
