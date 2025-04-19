package service

import (
	"context"
	"strings"

	"github.com/jinzhu/copier"
	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/model"
	"github.com/sztu/mutli-table/pkg"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/pkg/snowflake"
)

func RegisterSerivce(ctx context.Context, dto *DTO.SignUpRequestDTO) *apiError.ApiError {
	dto.Password = pkg.EncryptPassword(dto.Password)
	var user model.User

	err := copier.Copy(&user, dto)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}

	user.UserID, err = snowflake.GetID()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}

	err = dao.CreateUser(ctx, &user)
	if err != nil {
		// 检查是否为唯一约束冲突
		if err.Error() != "" && ( // 这里可以根据实际数据库驱动调整
		// 你也可以用 errors.As 或 errors.Is 判断具体错误类型
		// 这里只做字符串包含判断
		contains(err.Error(), "Duplicate entry") ||
			contains(err.Error(), "UNIQUE constraint failed")) {
			return &apiError.ApiError{
				Code: code.InvalidParam,
				Msg:  "用户名或邮箱已被使用",
			}
		}
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}
	return nil
}

// contains 是一个简单的字符串包含判断函数
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
