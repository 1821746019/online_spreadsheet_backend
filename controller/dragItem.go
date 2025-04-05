package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

func CreateDragCellHandler(c *gin.Context) {
	var req DTO.CreateDragItemRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("CreateDragCellHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()
	resp, apiErr := service.CreateDragItem(ctx, currentUserID, &req)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("待拖动单元格创建失败", zap.Any("error", apiErr))
		return
	}
	ResponseSuccess(c, resp)
}

func ListDragCellsHandler(c *gin.Context) {
	classID, err := strconv.ParseInt(c.Param("class_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的表格ID")
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid page")
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid page_size")
		return
	}
	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()
	resp, apiErr := service.ListDragItems(ctx, currentUserID, classID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("查询失败", zap.Any("error", apiErr))
		return
	}

	ResponseSuccess(c, resp)
}

func GetDragCellHandler(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("drag_item_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的元素ID")
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()
	resp, apiErr := service.GetDragItem(ctx, currentUserID, itemID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("获取失败", zap.Any("error", apiErr))
		return
	}

	ResponseSuccess(c, resp)
}

func UpdateDragCellHandler(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("drag_item_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的元素ID")
		return
	}

	var req DTO.UpdateDragItemRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("参数绑定失败", zap.Error(err))
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()

	resp, apiErr := service.UpdateDragItem(ctx, currentUserID, itemID, &req)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("更新失败", zap.Any("error", apiErr))
		return
	}

	ResponseSuccess(c, resp)
}

func DeleteDragCellHandler(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("drag_item_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的元素ID")
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()
	if apiErr := service.DeleteDragItem(ctx, currentUserID, itemID); apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("删除失败", zap.Any("error", apiErr))
		return
	}

	ResponseSuccess(c, nil)
}

func MoveDragItemHandler(c *gin.Context) {
	sheetIDStr := c.Param("sheet_id")
	dragItemIDStr := c.Param("drag_item_id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的表格ID")
		return
	}
	dragItemID, err := strconv.ParseInt(dragItemIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "无效的元素ID")
		return
	}
	var req DTO.MoveDragItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("MoveDragItemHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}
	ctx := c.Request.Context()
	err = service.MoveDragItem(ctx, currentUserID, sheetID, dragItemID, &req)
	if err != nil {
		ResponseErrorWithMsg(c, code.ServerError, "拖拽失败")
		return
	}
}
