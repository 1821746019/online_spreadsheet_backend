package controller

// // CreateDocumentHandler 创建文档接口
// // @Summary 创建文档接口
// // @Description 创建新的在线编辑文档
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Param document body DTO.CreateDocumentRequestDTO true "创建文档请求参数"
// // @Success 200 {object} Response
// // @Router /api/v1/document [post]
// func CreateDocumentHandler(c *gin.Context) {
// 	var req DTO.CreateDocumentRequestDTO
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
// 		zap.L().Error("CreateDocumentHandler.ShouldBindJSON() 失败", zap.Error(err))
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	resp, apiErr := service.CreateDocumentService(ctx, &req)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.CreateDocumentService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, resp)
// }

// // ListDocumentsHandler 查询文档列表接口
// // @Summary 查询文档列表接口
// // @Description 查询当前用户可访问的文档列表
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Success 200 {object} Response
// // @Router /api/v1/documents [get]
// func ListDocumentsHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	// 假设通过 JWT 或中间件从上下文中获取当前用户ID
// 	userID := GetUserIDFromContext(c) // 需要你实现此函数
// 	resp, apiErr := service.ListDocumentsService(ctx, userID)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.ListDocumentsService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, resp)
// }

// // GetDocumentHandler 获取文档详情接口
// // @Summary 获取文档详情接口
// // @Description 根据文档ID获取文档详情
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Param id path int true "文档ID"
// // @Success 200 {object} Response
// // @Router /api/v1/document/{id} [get]
// func GetDocumentHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, "无效的文档ID")
// 		zap.L().Error("GetDocumentHandler.ParseInt() 失败", zap.Error(err))
// 		return
// 	}

// 	resp, apiErr := service.GetDocumentService(ctx, id)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.GetDocumentService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, resp)
// }

// // UpdateDocumentHandler 更新文档接口
// // @Summary 更新文档接口
// // @Description 根据文档ID更新文档信息
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Param id path int true "文档ID"
// // @Param document body DTO.UpdateDocumentRequestDTO true "更新文档请求参数"
// // @Success 200 {object} Response
// // @Router /api/v1/document/{id} [put]
// func UpdateDocumentHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, "无效的文档ID")
// 		zap.L().Error("UpdateDocumentHandler.ParseInt() 失败", zap.Error(err))
// 		return
// 	}

// 	var req DTO.UpdateDocumentRequestDTO
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
// 		zap.L().Error("UpdateDocumentHandler.ShouldBindJSON() 失败", zap.Error(err))
// 		return
// 	}

// 	resp, apiErr := service.UpdateDocumentService(ctx, id, &req)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.UpdateDocumentService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, resp)
// }

// // DeleteDocumentHandler 删除文档接口（逻辑删除）
// // @Summary 删除文档接口
// // @Description 根据文档ID逻辑删除文档
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Param id path int true "文档ID"
// // @Success 200 {object} Response
// // @Router /api/v1/document/{id} [delete]
// func DeleteDocumentHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, "无效的文档ID")
// 		zap.L().Error("DeleteDocumentHandler.ParseInt() 失败", zap.Error(err))
// 		return
// 	}

// 	apiErr := service.DeleteDocumentService(ctx, id)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.DeleteDocumentService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, "删除成功")
// }

// // ShareDocumentHandler 文档共享接口
// // @Summary 文档共享接口
// // @Description 为文档设置共享权限
// // @Tags 文档
// // @Accept json
// // @Produce json
// // @Param id path int true "文档ID"
// // @Param share body DTO.ShareDocumentRequestDTO true "共享请求参数"
// // @Success 200 {object} Response
// // @Router /api/v1/document/{id}/share [post]
// func ShareDocumentHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, "无效的文档ID")
// 		zap.L().Error("ShareDocumentHandler.ParseInt() 失败", zap.Error(err))
// 		return
// 	}

// 	var req DTO.ShareDocumentRequestDTO
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
// 		zap.L().Error("ShareDocumentHandler.ShouldBindJSON() 失败", zap.Error(err))
// 		return
// 	}

// 	apiErr := service.ShareDocumentService(ctx, id, &req)
// 	if apiErr != nil {
// 		ResponseErrorWithApiError(c, apiErr)
// 		zap.L().Error("service.ShareDocumentService() 失败", zap.Error(apiErr))
// 		return
// 	}
// 	ResponseSuccess(c, "共享成功")
// }
