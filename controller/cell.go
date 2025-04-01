package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

// 获取单元格
func GetCellsHandler(c *gin.Context) {
	sheetID, _ := strconv.ParseInt(c.Param("sheet_id"), 10, 64)
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
	cells, apiError := service.GetCells(ctx, currentUserID, sheetID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("GetCells 失败", zap.Error(apiError))
		return
	}

	ResponseSuccess(c, cells)
}

// 更新单元格
func UpdateCellHandler(c *gin.Context) {
	sheetID, _ := strconv.ParseInt(c.Param("sheet_id"), 10, 64)

	var req DTO.UpdateCellRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("UpdateCellHandler.ShouldBindJSON() 失败", zap.Error(err))
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
	if err := service.UpdateCell(ctx, currentUserID, sheetID, req); err != nil {
		ResponseErrorWithApiError(c, err)
		zap.L().Error("UpdateCell 失败", zap.Error(err))
		return
	}

	ResponseSuccess(c, "更新成功")
}
