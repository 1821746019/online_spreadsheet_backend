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

// CreateSheet 创建新的工作表，同时为创建者添加 ADMIN 权限
func CreateSheet(ctx context.Context, userID, classID int64, dto *DTO.CreateSheetRequestDTO) (*DTO.SheetResponseDTO, *apiError.ApiError) {
	_, err := dao.GetClassByID(ctx, classID)
	if err != nil {
		zap.L().Error("CreateSheet 查询班级失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "创建工作表失败"}
	}
	// 获取数据库句柄，并开启事务
	db := mysql.GetDB().WithContext(ctx)
	tx := db.Begin()
	// 保证事务回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 构造 Sheet 对象
	sheet := &model.Sheet{
		Name:       dto.Name,
		CreatorID:  userID,
		Week:       int32(dto.Week),
		Row:        int32(dto.Row),
		Col:        int32(dto.Col),
		ClassID:    classID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 插入 Sheet 记录
	if err := dao.CreateSheetTx(context.Background(), tx, sheet); err != nil {
		tx.Rollback()
		zap.L().Error("CreateSheet 失败：插入 sheet 记录错误", zap.Error(err))
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "创建工作表失败",
		}
	}

	var cells []model.Cell
	for row := 1; row <= dto.Row; row++ {
		for col := 1; col <= dto.Col; col++ {
			cells = append(cells, model.Cell{
				SheetID:    sheet.ID,
				RowIndex:   int32(row),
				ColIndex:   int32(col),
				ItemID:     nil,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
		}
	}

	// 批量插入 Cells
	if err := dao.CreateBatchCellsTx(tx, ctx, cells); err != nil {
		tx.Rollback()
		zap.L().Error("CreateSheet 失败：批量插入 cell 记录错误", zap.Error(err))
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "初始化单元格失败",
		}
	}

	// 构造 Permission 记录
	// perm := &model.Permission{
	// 	UserID:     sheet.CreatorID,
	// 	SheetID:    sheet.ID,
	// 	CreateTime: time.Now(),
	// 	UpdateTime: time.Now(),
	// }

	// // 插入 Permission 记录
	// if err := dao.CreatePermissionTx(tx, perm); err != nil {
	// 	tx.Rollback()
	// 	zap.L().Error("CreateSheet 失败：插入 permission 记录错误", zap.Error(err))
	// 	return nil, &apiError.ApiError{
	// 		Code: code.ServerError,
	// 		Msg:  "创建权限记录失败",
	// 	}
	// }

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.L().Error("CreateSheet 失败：事务提交错误", zap.Error(err))
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "事务提交失败",
		}
	}

	// 返回成功的响应 DTO
	return &DTO.SheetResponseDTO{
		ID:         sheet.ID,
		Name:       sheet.Name,
		CreatorID:  sheet.CreatorID,
		Week:       int(sheet.Week),
		Row:        dto.Row,
		Col:        dto.Col,
		ClassID:    classID,
		CreateTime: time.Now().String(),
		UpdateTime: time.Now().String(),
	}, nil
}

// ListSheets 获取所有的工作表列表
func ListSheets(ctx context.Context, userID, classID int64, page, pageSize int) (*DTO.SheetListResponseDTO, *apiError.ApiError) {
	_, err := dao.GetClassByID(ctx, classID)
	if err != nil {
		zap.L().Error("CreateSheet 查询班级失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "创建工作表失败"}
	}
	sheets, total, err := dao.ListSheets(ctx, userID, classID, page, pageSize)
	if err != nil {
		zap.L().Error("ListSheets 查询失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "获取工作表列表失败"}
	}

	var sheetDTOs []DTO.SheetDTO
	for _, s := range sheets {
		sheetDTOs = append(sheetDTOs, DTO.SheetDTO{
			ID:      s.ID,
			Name:    s.Name,
			ClassID: s.ClassID,
		})
	}

	return &DTO.SheetListResponseDTO{
		Total:  total,
		Page:   page,
		Sheets: sheetDTOs,
	}, nil
}

// GetSheet 根据 sheetID 获取工作表详情
func GetSheet(ctx context.Context, userID int64, sheetID int64) (*DTO.SheetDetailResponseDTO, *apiError.ApiError) {
	// 查询工作表详情
	sheet, err := dao.GetSheetByID(ctx, sheetID)
	if err != nil {
		zap.L().Error("GetSheet 查询工作表失败", zap.Error(err))
		return nil, &apiError.ApiError{Code: code.ServerError, Msg: "查询工作表失败"}
	}
	if sheet == nil {
		return nil, &apiError.ApiError{Code: code.NotFound, Msg: "工作表不存在"}
	}

	// 构造返回的 DTO
	return &DTO.SheetDetailResponseDTO{
		ID:   sheet.ID,
		Name: sheet.Name,
		Week: int(sheet.Week),
		Row:  int(sheet.Row),
		Col:  int(sheet.Col),
	}, nil
}

// UpdateSheet 更新工作表信息
func UpdateSheet(ctx context.Context, userID, sheetID int64, dto *DTO.UpdateSheetRequestDTO) *apiError.ApiError {
	sheet, err := dao.GetSheetByID(ctx, sheetID)
	if err != nil {
		zap.L().Error("UpdateSheet 查询失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "更新工作表失败"}
	}
	if sheet == nil {
		return &apiError.ApiError{Code: code.NotFound, Msg: "工作表不存在"}
	}

	// 根据传入非 nil 的字段更新工作表
	if dto.Name != nil {
		sheet.Name = *dto.Name
	}
	sheet.UpdateTime = time.Now()

	if err := dao.UpdateSheet(ctx, sheet); err != nil {
		zap.L().Error("UpdateSheet 更新失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "更新工作表失败"}
	}
	return nil
}

// DeleteSheet 逻辑删除工作表
func DeleteSheet(ctx context.Context, sheetID int64) *apiError.ApiError {
	if err := dao.DeleteSheet(ctx, sheetID); err != nil {
		zap.L().Error("DeleteSheet 删除失败", zap.Error(err))
		return &apiError.ApiError{Code: code.ServerError, Msg: "删除工作表失败"}
	}
	return nil
}
