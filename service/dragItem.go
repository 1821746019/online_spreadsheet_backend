package service

import (
	"context"
	"time"

	dao "github.com/sztu/mutli-table/DAO"
	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/model"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"go.uber.org/zap"
)

func CreateDragItem(ctx context.Context, userID int64, sheetID int64, req *DTO.CreateDragItemRequestDTO) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("CreateDragItem 查询权限失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}

	item := &model.DragItem{
		Content:   req.Content,
		CreatorID: userID,

		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := dao.CreateDraggableItem(ctx, item); err != nil {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "创建失败"}
	}

	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		Content:    item.Content,
		SheetID:    sheetID,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func ListDragItems(ctx context.Context, userID int64, sheetID int64) ([]*DTO.DragItemResponseDTO, *apiError.ApiError) {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("ListDragItems 查询权限失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}

	// 获取该工作表所有的拖拽元素
	items, err := dao.ListDraggableItemsBySheet(ctx, sheetID)
	if err != nil {
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询拖拽元素失败"}
	}

	// 获取该工作表所有的单元格数据
	cells, err := dao.GetCellsBySheetID(ctx, sheetID)
	if err != nil {
		zap.L().Error("ListDragItems 查询单元格失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询单元格失败"}
	}
	// 建立一个map，用来记录已经被单元格关联的拖拽元素ID
	usedItems := make(map[int64]bool)
	for _, cell := range cells {
		// 如果cell的item_id不为0（或非nil，根据具体实现判断），则说明已经被关联
		if cell.ItemID != nil {
			usedItems[*cell.ItemID] = true
		}
	}

	// 过滤掉已经关联的元素，只返回未被使用的拖拽元素
	res := make([]*DTO.DragItemResponseDTO, 0, len(items))
	for _, item := range items {
		if usedItems[item.ID] {
			continue
		}
		res = append(res, &DTO.DragItemResponseDTO{
			ID:         item.ID,
			SheetID:    sheetID,
			Content:    item.Content,
			CreatorID:  item.CreatorID,
			CreateTime: item.CreateTime.Format(time.RFC3339),
			UpdateTime: item.UpdateTime.Format(time.RFC3339),
		})
	}
	return res, nil
}

func GetDragItem(ctx context.Context, userID int64, sheetID int64, itemID int64) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("GetDragItem 查询权限失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, &apiError.ApiError{Code: code.NotFound, Msg: "元素不存在"}
	}

	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		SheetID:    item.SheetID,
		Content:    item.Content,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func UpdateDragItem(ctx context.Context, userID int64, sheetID int64, itemID int64, req *DTO.UpdateDragItemRequestDTO) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("GetDragItem 查询权限失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, &apiError.ApiError{code.NotFound, "元素不存在"}
	}
	item.Content = req.Content
	item.UpdateTime = time.Now()

	if err := dao.UpdateDraggableItem(ctx, item); err != nil {
		return nil, &apiError.ApiError{code.ServerError, "更新失败"}
	}

	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		SheetID:    item.SheetID,
		Content:    item.Content,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func DeleteDragItem(ctx context.Context, userID int64, sheetID int64, itemID int64) *apiError.ApiError {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("GetDragItem 查询权限失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil {
		zap.L().Error("DeleteDragItem 查询元素失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{code.ServerError, "查询元素失败"}
	}
	if item == nil {
		return &apiError.ApiError{code.NotFound, "元素不存在"}
	}
	// 5. 检查元素是否被单元格引用
	var refCount int64
	if err := mysql.GetDB().WithContext(ctx).
		Model(&model.Cell{}).
		Where("item_id = ? AND delete_time = 0", itemID).
		Count(&refCount).Error; err != nil {
		zap.L().Error("DeleteDragItem 检查引用失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{code.ServerError, "系统繁忙，请稍后再试"}
	}
	if refCount > 0 {
		return &apiError.ApiError{code.ServerError, "存在关联单元格，请先解除关联"}
	}
	// 6. 执行删除操作
	if err := dao.DeleteDraggableItem(ctx, itemID, sheetID); err != nil {
		zap.L().Error("DeleteDragItem 删除失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{code.ServerError, "删除操作失败"}
	}
	return nil
}

// MoveDragItem 实现拖拽元素的移动或交换
// 业务逻辑：
// 1. 如果拖拽元素原本在某个单元格中（sourceCell 不为 nil）
//   - 当目标单元格已有拖拽元素时，交换两个单元格的 item_id
//   - 当目标单元格为空时，直接移动拖拽元素（同时清空原单元格的关联）
//
// 2. 如果拖拽元素原本不在任何单元格中（sourceCell 为 nil，即在待拖拽列表中）
//   - 当目标单元格已有拖拽元素时，返回错误（或根据业务设计可定义其他处理）
//   - 当目标单元格为空时，从待拖拽列表中获取该拖拽元素，将目标单元格的 Content 设为拖拽元素的 Content，并关联该拖拽元素
func MoveDragItem(ctx context.Context, userID, sheetID, dragItemID int64, dto *DTO.MoveDragItemRequest) *apiError.ApiError {
	perm, err := dao.GetPermission(ctx, userID, sheetID)
	if err != nil {
		zap.L().Error("GetDragItem 查询权限失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "检查权限失败"}
	}
	if perm == nil || perm.AccessLevel == "READ" {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限修改该工作表"}
	}
	db := mysql.GetDB().WithContext(ctx)
	tx := db.Begin()
	// 保证事务异常或 panic 时回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取拖拽元素当前所在单元格（可能为 nil，代表该拖拽元素在待拖拽列表中）
	sourceCell, err := dao.GetCellByDragItemIDTx(ctx, tx, sheetID, dragItemID)
	if err != nil {
		tx.Rollback()
		zap.L().Error("MoveDragItem 获取源单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取源单元格失败"}
	}

	// 获取目标单元格（必须存在，否则请确保工作表创建时已生成所有单元格）
	targetCell, err := dao.GetCellByPositionTx(ctx, tx, sheetID, dto.TargetRow, dto.TargetCol)
	if err != nil {
		tx.Rollback()
		zap.L().Error("MoveDragItem 获取目标单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取目标单元格失败"}
	}

	// 如果目标单元格已有拖拽元素（targetCell.ItemID != nil）
	if targetCell.ItemID != nil {
		// 若拖拽元素原本在单元格中，则执行交换操作
		if sourceCell != nil {
			tmp := targetCell.ItemID
			targetCell.ItemID = sourceCell.ItemID
			sourceCell.ItemID = tmp

			// 更新两个单元格
			if err := dao.UpdateCellTx(ctx, tx, sourceCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新源单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
			if err := dao.UpdateCellTx(ctx, tx, targetCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新目标单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
		} else {
			// 若拖拽元素来自待拖拽列表，但目标单元格已有拖拽元素，则按业务规则不支持交换
			tx.Rollback()
			zap.L().Error("MoveDragItem 错误：待拖拽列表中的元素不能与已有元素交换", zap.Int64("dragItemID", dragItemID))
			return &apiError.ApiError{Code: code.InvalidParam, Msg: "目标单元格已有拖拽元素，无法进行交换"}
		}
	} else {
		// 目标单元格没有拖拽元素
		if sourceCell == nil {
			// 拖拽元素来自待拖拽列表，从待拖拽库中获取该拖拽元素记录
			dragItem, err := dao.GetDraggableItemByIDTx(ctx, tx, sheetID, dragItemID)
			if err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 获取待拖拽元素失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "获取待拖拽元素失败"}
			}
			if dragItem == nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 未找到待拖拽元素", zap.Int64("dragItemID", dragItemID))
				return &apiError.ApiError{Code: code.NotFound, Msg: "拖拽元素不存在"}
			}
			// 将目标单元格的内容设为拖拽元素的内容，并关联该拖拽元素
			targetCell.Content = dragItem.Content
			targetCell.ItemID = &dragItemID
			if err := dao.UpdateCellTx(ctx, tx, targetCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新目标单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
		} else {
			// 拖拽元素原本在某单元格中，且目标单元格为空：直接移动
			targetCell.Content = "" // 清空目标原有内容
			targetCell.ItemID = &dragItemID
			if err := dao.UpdateCellTx(ctx, tx, targetCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新目标单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
			// 清空原单元格（如果与目标单元格不在同一位置）
			if sourceCell.RowIndex != targetCell.RowIndex || sourceCell.ColIndex != targetCell.ColIndex {
				sourceCell.ItemID = nil
				if err := dao.UpdateCellTx(ctx, tx, sourceCell); err != nil {
					tx.Rollback()
					zap.L().Error("MoveDragItem 清空源单元格失败", zap.Error(err))
					return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
				}
			}
		}
	}

	// 更新单元格更新时间
	targetCell.UpdateTime = time.Now()
	if sourceCell != nil {
		sourceCell.UpdateTime = time.Now()
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.L().Error("MoveDragItem 事务提交失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "事务提交失败"}
	}

	return nil
}
