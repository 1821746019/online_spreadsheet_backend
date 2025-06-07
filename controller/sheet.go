package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

func CreateSheetHandler(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
		return
	}
	var req DTO.CreateSheetRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("CreateSheetHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	// 从上下文中获取当前用户ID
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
	resp, apiError := service.CreateSheet(ctx, currentUserID, classID, &req)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("CreateSheet 失败", zap.Error(apiError))
		return
	}

	ResponseSuccess(c, resp)
}

// ListSheetsHandler 获取工作表列表（支持分页查询）
func ListSheetsHandler(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
		return
	}
	// 默认分页参数：第 1 页，每页 10 条记录
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
	result, apiErr := service.ListSheets(ctx, currentUserID, classID, page, pageSize)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("ListSheetsHandler failed", zap.Error(apiErr))
		return
	}
	ResponseSuccess(c, result)
}

// GetSheetHandler 获取单个工作表详情
func GetSheetHandler(c *gin.Context) {
	sheetIDStr := c.Param("sheet_id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
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
	sheet, apiErr := service.GetSheet(ctx, currentUserID, sheetID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("GetSheetHandler failed", zap.Error(apiErr))
		return
	}
	ResponseSuccess(c, sheet)
}

// UpdateSheetHandler 更新工作表信息
func UpdateSheetHandler(c *gin.Context) {
	sheetIDStr := c.Param("sheet_id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
		return
	}

	var req DTO.UpdateSheetRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("UpdateSheetHandler binding 失败", zap.Error(err))
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
	apiErr := service.UpdateSheet(ctx, currentUserID, sheetID, &req)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("UpdateSheetHandler 失败", zap.Error(apiErr))
		return
	}
	ResponseSuccess(c, "更新成功")
}

// DeleteSheetHandler 删除工作表（逻辑删除）
func DeleteSheetHandler(c *gin.Context) {
	sheetIDStr := c.Param("sheet_id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
		return
	}
	ctx := c.Request.Context()
	apiErr := service.DeleteSheet(ctx, sheetID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("DeleteSheetHandler failed", zap.Error(apiErr))
		return
	}
	ResponseSuccess(c, "删除成功")
}
