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
	if err != nil {
		zap.L().Error("GetCellByRowAndCol 查询单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取单元格失败"}
	}
	if targetCell.ItemID == nil {
		return &apiError.ApiError{Code: code.InvalidParam, Msg: "删除的单元格为空"}
	}
	dragItem, err := dao.GetDraggableItemByID(ctx, *targetCell.ItemID)
	if err != nil {
		zap.L().Error("GetDraggableItemByID 查询拖拽项失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取拖拽项失败"}
	}
	if dragItem.CreatorID != userID {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限删除该拖拽项"}
	}
	targetCell.ItemID = nil
	targetCell.UpdateTime = time.Now()
	targetCell.LastModifiedBy = userID
	currentSheet, err := dao.GetSheetByID(ctx, sheetID)
	if err != nil || currentSheet == nil {
		zap.L().Error("获取工作表信息失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取工作表失败"}
	}

	totalWeeks, err := dao.GetClassTotalWeeks(ctx, currentSheet.ClassID)
	if err != nil {
		zap.L().Error("获取班级总周数失败",
			zap.Int64("classID", currentSheet.ClassID),
			zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取班级周数失败"}
	}

	// 根据周类型生成目标周列表
	var targetWeeks []int
	switch dragItem.WeekType {
	case "single":
		for w := 1; w <= totalWeeks; w += 2 { // 假设总共有18周
			targetWeeks = append(targetWeeks, w)
		}
	case "double":
		for w := 2; w <= totalWeeks; w += 2 {
			targetWeeks = append(targetWeeks, w)
		}
	case "all":
		for w := 1; w <= totalWeeks; w++ {
			targetWeeks = append(targetWeeks, w)
		}
	}
	// 为每个目标周创建/更新单元格
	for _, week := range targetWeeks {
		if week == int(currentSheet.Week) { // 跳过当前周（已处理）
			continue
		}

		// 获取目标周的工作表
		targetSheet, err := dao.GetSheetByClassIDandWeek(ctx, currentSheet.ClassID, week)
		if err != nil || targetSheet == nil {
			zap.L().Error("获取周工作表失败",
				zap.Int("week", week),
				zap.Error(err))
			continue
		}
		// 获取目标单元格
		targetCell, err := dao.GetCellByPosition(ctx, targetSheet.ID, req.Row, req.Col)
		if err != nil || targetCell == nil {
			zap.L().Error("获取目标单元格失败",
				zap.Int("week", week),
				zap.Error(err))
			continue
		}

		targetCell.ItemID = nil
		targetCell.UpdateTime = time.Now()
		targetCell.LastModifiedBy = userID
		// 更新目标单元格
		if err := dao.UpdateCell(ctx, targetSheet.ID, targetCell); err != nil {
			zap.L().Error("更新周单元格失败",
				zap.Int("week", week),
				zap.Error(err))
		}
	}

	return &apiError.ApiError{Code: code.Success, Msg: "操作成功"}
}
