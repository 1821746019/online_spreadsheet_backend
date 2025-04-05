package service

import (
	"context"
	"time"

	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/model"
	"github.com/sztu/mutli-table/pkg/apiError"
	"github.com/sztu/mutli-table/pkg/code"
	"go.uber.org/zap"
)

func CreateClass(ctx context.Context, userID int64, req *DTO.CreateClassRequestDTO) (*DTO.ClassResponseDTO, *apiError.ApiError) {
	if exist, _ := dao.ClassNameExists(ctx, req.Name); exist {
		return nil, &apiError.ApiError{Code: code.InvalidParam, Msg: "班级名称已存在"}
	}
	// 创建班级
	class := &model.Class{
		Name:       req.Name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := dao.CreateClass(ctx, class); err != nil {
		zap.L().Error("创建班级失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "创建班级失败"}
	}
	// 返回成功响应
	return &DTO.ClassResponseDTO{ID: class.ID, Name: class.Name}, nil
}

func ListClasses(ctx context.Context, userID int64, page, pageSize int) (*DTO.ClassListDTO, *apiError.ApiError) {
	classes, total, err := dao.ListClasses(ctx, userID, page, pageSize)
	if err != nil {
		zap.L().Error("查询班级列表失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询班级列表失败"}
	}
	// 构建响应
	classList := make([]DTO.ClassSimpleItemDTO, 0, len(classes))
	for _, class := range classes {
		classList = append(classList, DTO.ClassSimpleItemDTO{ID: class.ID, Name: class.Name})
	}
	return &DTO.ClassListDTO{Total: total, List: classList}, nil
}

func GetClass(ctx context.Context, userID, classID int64) (*DTO.ClassResponseDTO, *apiError.ApiError) {
	class, err := dao.GetClassByID(ctx, classID)
	if err != nil {
		zap.L().Error("查询班级失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询班级失败"}
	}
	if class == nil {
		return nil, &apiError.ApiError{Code: code.NotFound, Msg: "班级不存在"}
	}
	return &DTO.ClassResponseDTO{ID: class.ID, Name: class.Name}, nil
}

func UpdateClass(ctx context.Context, userID, classID int64, req *DTO.UpdateClassRequestDTO) *apiError.ApiError {
	class, err := dao.GetClassByID(ctx, classID)
	if err != nil {
		zap.L().Error("查询班级失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "查询班级失败"}
	}
	if class == nil {
		return &apiError.ApiError{Code: code.NotFound, Msg: "班级不存在"}
	}
	// 更新班级信息
	class.Name = req.Name
	class.UpdateTime = time.Now()

	if err := dao.UpdateClass(ctx, class); err != nil {
		zap.L().Error("更新班级失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "更新班级失败"}
	}
	return nil
}

func DeleteClass(ctx context.Context, userID, classID int64) *apiError.ApiError {
	class, err := dao.GetClassByID(ctx, classID)
	if err != nil {
		zap.L().Error("查询班级失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "查询班级失败"}
	}
	if class == nil {
		return &apiError.ApiError{Code: code.NotFound, Msg: "班级不存在"}
	}
	if err := dao.DeleteClass(ctx, classID); err != nil {
		zap.L().Error("删除班级失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "删除班级失败"}
	}
	return nil
}
