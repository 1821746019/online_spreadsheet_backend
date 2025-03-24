package controller

import (
	"github.com/sztu/mutli-table/pkg/code"

	"github.com/gin-gonic/gin"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

// SignUpHandler 注册接口
// @Summary 注册接口
// @Description 注册接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Param email body string true "邮箱"
// @Success 200 {object} Response
// @Router /api/v1/signup [post]
func SignUpHandler(c *gin.Context) {
	var SignupDTO DTO.SignUpRequestDTO
	if err := c.ShouldBindJSON(&SignupDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("SignUpHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	ctx := c.Request.Context()

	apiError := service.RegisterSerivce(ctx, &SignupDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("AuthServiceInterface.SignupService() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, nil)
}

// LoginHandler 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} Response
// @Router /api/v1/login [post]
func LoginHandler(c *gin.Context) {
	var loginDTO DTO.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("LoginHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}

	ctx := c.Request.Context()

	resp, apiError := service.LoginService(ctx, &loginDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("AuthServiceInterface.LoginService() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, resp)
	return
}

// LogoutHandler 退出登录
// @Summary 退出登录
// @Description 退出登录
// @Tags 登录
// @Accept json
// @Produce json
// @Param access_token query string true
// @Param refresh_token query string true
// @Success 200 {object} Response
// @Router /api/v1/logout [post]
func LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accessToken := c.Query("access_token")
	apiError := service.LogoutService(ctx, accessToken)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
	return
}
