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

func CreateDragItem(ctx context.Context, userID int64, req *DTO.CreateDragItemRequestDTO) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	// 复用逻辑验证
	if len(req.SelectedClassIDs) == 0 {
		return nil, &apiError.ApiError{Code: code.InvalidParam, Msg: "请选择要关联的班级"}
	}

	item := &model.DraggableItem{
		Content:    req.Content,
		CreatorID:  userID,
		WeekType:   req.WeekType,
		Classroom:  req.Classroom,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 使用事务处理创建操作
	tx := mysql.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := dao.CreateDraggableItemTx(ctx, tx, item); err != nil {
		tx.Rollback()
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "创建失败"}
	}

	// 关联多个班级
	for _, classID := range req.SelectedClassIDs {
		if err := dao.CreateItemSheetRelationTx(ctx, tx, item.ID, classID); err != nil {
			tx.Rollback()
			zap.L().Error("创建班级关联失败",
				zap.Int64("classID", classID),
				zap.Error(err))
			return nil, &apiError.ApiError{Code: code.ServerError, Msg: "关联班级失败"}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "事务提交失败"}
	}

	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		WeekType:   item.WeekType,
		Classroom:  item.Classroom,
		Content:    item.Content,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func ListDragItems(ctx context.Context, userID int64, classID int64) ([]*DTO.DragItemResponseDTO, *apiError.ApiError) {
	// 班级存在性校验
	if _, err := dao.GetClassByID(ctx, classID); err != nil {
		return nil, &apiError.ApiError{Code: code.NotFound, Msg: "班级不存在"}
	}
	// 获取班级关联的所有元素
	items, err := dao.ListDraggableItemsByClass(ctx, classID)
	if err != nil {
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询拖拽元素失败"}
	}
	res := make([]*DTO.DragItemResponseDTO, 0, len(items))
	for _, item := range items {
		res = append(res, &DTO.DragItemResponseDTO{
			ID:         item.ID,
			Content:    item.Content,
			WeekType:   item.WeekType,
			Classroom:  item.Classroom,
			CreatorID:  item.CreatorID,
			CreateTime: item.CreateTime.Format(time.RFC3339),
			UpdateTime: item.UpdateTime.Format(time.RFC3339),
		})
	}
	return res, nil
}

func GetDragItem(ctx context.Context, userID int64, itemID int64) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, &apiError.ApiError{Code: code.NotFound, Msg: "元素不存在"}
	}
	if item.CreatorID != userID {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限读取该元素"}
	}
	classNames, err := dao.GetClassNamesByItemID(ctx, itemID)
	if err != nil {
		zap.L().Error("获取班级名称失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "获取班级信息失败"}
	}
	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		Content:    item.Content,
		WeekType:   item.WeekType,
		Classroom:  item.Classroom,
		ClassNames: classNames,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func UpdateDragItem(ctx context.Context, userID int64, itemID int64, req *DTO.UpdateDragItemRequestDTO) (*DTO.DragItemResponseDTO, *apiError.ApiError) {
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, &apiError.ApiError{code.NotFound, "元素不存在"}
	}
	if item.CreatorID != userID {
		return nil, &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限读取该元素"}
	}
	// 开启事务
	tx := mysql.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新基础信息
	item.Content = req.Content
	item.UpdateTime = time.Now()
	item.WeekType = req.WeekType
	item.Classroom = req.Classroom
	if err := dao.UpdateDraggableItemTx(ctx, tx, item); err != nil {
		tx.Rollback()
		return nil, &apiError.ApiError{code.ServerError, "基础信息更新失败"}
	}

	// 删除旧的班级关联
	if err := dao.DeleteItemClassRelationsTx(ctx, tx, itemID); err != nil {
		tx.Rollback()
		zap.L().Error("删除旧班级关联失败", zap.Int64("itemID", itemID), zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "班级关联更新失败"}
	}

	// 添加新的班级关联
	for _, classID := range req.SelectedClassIDs {
		if err := dao.CreateItemSheetRelationTx(ctx, tx, itemID, classID); err != nil {
			tx.Rollback()
			zap.L().Error("创建班级关联失败",
				zap.Int64("classID", classID),
				zap.Error(err))
			return nil, &apiError.ApiError{Code: code.ServerError, Msg: "班级关联更新失败"}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "事务提交失败"}
	}

	// 获取最新班级名称
	classNames, _ := dao.GetClassNamesByItemID(ctx, itemID) // 忽略错误，主流程已成功

	return &DTO.DragItemResponseDTO{
		ID:         item.ID,
		Content:    item.Content,
		WeekType:   item.WeekType,
		Classroom:  item.Classroom,
		ClassNames: classNames,
		CreatorID:  item.CreatorID,
		CreateTime: item.CreateTime.Format(time.RFC3339),
		UpdateTime: item.UpdateTime.Format(time.RFC3339),
	}, nil
}

func DeleteDragItem(ctx context.Context, userID int64, itemID int64) *apiError.ApiError {
	item, err := dao.GetDraggableItemByID(ctx, itemID)
	if err != nil || item == nil {
		return &apiError.ApiError{code.NotFound, "元素不存在"}
	}
	if item.CreatorID != userID {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限读取该元素"}
	}
	// 检查元素是否被单元格引用
	refCount, err := dao.CountCellReferences(ctx, itemID)
	if err != nil {
		zap.L().Error("DeleteDragItem 检查引用失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{code.ServerError, "系统繁忙，请稍后再试"}
	}
	if refCount > 0 {
		return &apiError.ApiError{code.ServerError, "存在关联单元格，请先解除关联"}
	}
	// 执行删除操作
	// 开启事务
	tx := mysql.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除班级关联关系
	if err := dao.DeleteItemClassRelationsTx(ctx, tx, itemID); err != nil {
		tx.Rollback()
		zap.L().Error("删除班级关联失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "删除关联信息失败"}
	}

	// 执行删除操作
	if err := dao.DeleteDraggableItemTx(ctx, tx, itemID); err != nil {
		tx.Rollback()
		zap.L().Error("DeleteDragItem 删除失败",
			zap.Int64("itemID", itemID),
			zap.Error(err))
		return &apiError.ApiError{code.ServerError, "删除操作失败"}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &apiError.ApiError{Code: code.ServerError, Msg: "事务提交失败"}
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
//   - 当目标单元格已有拖拽元素时，返回错误
//   - 当目标单元格为空时，从待拖拽列表中获取该拖拽元素，并关联该拖拽元素
func MoveDragItem(ctx context.Context, userID, sheetID, dragItemID int64, dto *DTO.MoveDragItemRequest) *apiError.ApiError {
	item, err := dao.GetDraggableItemByID(ctx, dragItemID)
	if err != nil || item == nil {
		return &apiError.ApiError{code.NotFound, "元素不存在"}
	}
	if item.CreatorID != userID {
		return &apiError.ApiError{Code: code.NoPermission, Msg: "没有权限读取该元素"}
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

	// 获取目标单元格
	targetCell, err := dao.GetCellByPositionTx(ctx, tx, sheetID, dto.TargetRow, dto.TargetCol)
	if err != nil || targetCell == nil {
		tx.Rollback()
		zap.L().Error("MoveDragItem 获取目标单元格失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取目标单元格失败"}
	}
	if targetCell.ItemID != nil && targetCell.LastModifiedBy != userID {
		// 目标单元格已有拖拽元素，不允许移动
		tx.Rollback()
		return &apiError.ApiError{Code: code.ServerError, Msg: "目标单元格已有拖拽元素"}
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
			targetCell.ItemID = &item.ID // 保持目标单元格的拖拽元素不变，或者根据业务逻辑调整
			if err := dao.UpdateCellTx(ctx, tx, targetCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新源单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
		}
	} else {
		// 目标单元格没有拖拽元素
		if sourceCell == nil {
			// 拖拽元素来自待拖拽列表，从待拖拽库中获取该拖拽元素记录
			dragItem, err := dao.GetDraggableItemByIDTx(ctx, tx, dragItemID)
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
			targetCell.ItemID = &dragItemID
			if err := dao.UpdateCellTx(ctx, tx, targetCell); err != nil {
				tx.Rollback()
				zap.L().Error("MoveDragItem 更新目标单元格失败", zap.Error(err))
				return &apiError.ApiError{Code: code.ServerError, Msg: "更新单元格失败"}
			}
		} else {
			// 拖拽元素原本在某单元格中，且目标单元格为空：直接移动
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

	currentSheet, err := dao.GetSheetByID(ctx, sheetID)
	if err != nil || currentSheet == nil {
		zap.L().Error("获取工作表信息失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "获取工作表失败"}
	}
	// 根据周类型生成目标周列表
	var targetWeeks []int
	switch item.WeekType {
	case "single":
		for w := 1; w <= 18; w += 2 { // 假设总共有18周
			targetWeeks = append(targetWeeks, w)
		}
	case "double":
		for w := 2; w <= 18; w += 2 {
			targetWeeks = append(targetWeeks, w)
		}
	case "all":
		for w := 1; w <= 18; w++ {
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
		if err != nil {
			zap.L().Error("获取周工作表失败",
				zap.Int("week", week),
				zap.Error(err))
			continue
		}
		// 获取目标单元格
		targetCell, err := dao.GetCellByPosition(ctx, targetSheet.ID, dto.TargetRow, dto.TargetCol)
		if err != nil || targetCell == nil {
			zap.L().Error("获取目标单元格失败",
				zap.Int("week", week),
				zap.Error(err))
			continue
		}

		targetCell.ItemID = &dragItemID
		targetCell.UpdateTime = time.Now()

		// 更新目标单元格
		if err := dao.UpdateCell(ctx, sheetID, targetCell); err != nil {
			zap.L().Error("更新周单元格失败",
				zap.Int("week", week),
				zap.Error(err))
		}
	}

	return nil
}
