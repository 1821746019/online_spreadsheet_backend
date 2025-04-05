package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

func CreateClassHandler(c *gin.Context) {
	var req DTO.CreateClassRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("CreateClassHandler.ShouldBindJSON() 失败", zap.Error(err))
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
	resp, apiErr := service.CreateClass(ctx, currentUserID, &req)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("CreateClass 失败", zap.Error(apiErr))
		return
	}

	ResponseSuccess(c, resp)
}

func ListClassesHandler(c *gin.Context) {
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
	result, apiErr := service.ListClasses(ctx, currentUserID, page, pageSize)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("ListClasses 失败", zap.Error(apiErr))
		return
	}
	ResponseSuccess(c, result)
}

func GetClassHandler(c *gin.Context) {
	classID, err := strconv.ParseInt(c.Param("class_id"), 10, 64)
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
	class, apiErr := service.GetClass(ctx, currentUserID, classID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("GetClass 失败", zap.Error(apiErr))
		return
	}

	ResponseSuccess(c, class)
}

func UpdateClassHandler(c *gin.Context) {
	classID, err := strconv.ParseInt(c.Param("class_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid class_id")
		return
	}

	var req DTO.UpdateClassRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("UpdateClassHandler.ShouldBindJSON() 失败", zap.Error(err))
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
	apiErr := service.UpdateClass(ctx, currentUserID, classID, &req)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("UpdateClass 失败", zap.Error(apiErr))
		return
	}

	ResponseSuccess(c, "更新成功")
}

func DeleteClassHandler(c *gin.Context) {
	classID, err := strconv.ParseInt(c.Param("class_id"), 10, 64)
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
	apiErr := service.DeleteClass(ctx, currentUserID, classID)
	if apiErr != nil {
		ResponseErrorWithApiError(c, apiErr)
		zap.L().Error("DeleteClass 失败", zap.Error(apiErr))
		return
	}

	ResponseSuccess(c, "删除成功")
}
