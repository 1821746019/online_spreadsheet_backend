package service

import (
	"context"
	"time"

	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/cache"
	"github.com/sztu/mutli-table/pkg"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/pkg/jwt"
)

// LoginService 登录服务
func LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError) {
	user, err := dao.FindUserByUsername(ctx, dto.Username)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "登录失败",
		}
	}
	if user == nil {
		return nil, &apiError.ApiError{
			Code: code.UserNotExist,
			Msg:  "用户不存在",
		}
	}
	if pkg.EncryptPassword(dto.Password) != user.Password {
		return nil, &apiError.ApiError{
			Code: code.PasswordError,
			Msg:  "密码错误",
		}
	}
	accessToken, err := jwt.GenerateToken(user.UserID, user.Username)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}

	return &DTO.LoginResponseDTO{
		AccessToken: accessToken,
		UserID:      user.UserID,
		Username:    user.Username,
	}, nil
}

func LogoutService(ctx context.Context, token ...string) *apiError.ApiError {
	for _, t := range token {
		myClaims, err := jwt.ParseToken(t)
		if err != nil {
			return &apiError.ApiError{
				Code: code.UserRefreshTokenError,
				Msg:  err.Error(),
			}
		}

		err = cache.AddTokenToBlacklist(ctx, t, time.Until(myClaims.ExpiresAt.Time))
		if err != nil {
			return &apiError.ApiError{
				Code: code.ServerError,
				Msg:  "登出失败",
			}
		}
	}

	return nil
}
