package DTO

type CreateSheetRequestDTO struct {
	Name string `json:"name" binding:"required"`
	Week int    `json:"week" binding:"required"`
	Row  int    `json:"row" binding:"required"`
	Col  int    `json:"col" binding:"required"`
}

type SheetResponseDTO struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreatorID  int64  `json:"creator_id"`
	Week       int    `json:"week"`
	Row        int    `json:"row"`
	Col        int    `json:"col"`
	ClassID    int64  `json:"class_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type UpdateSheetRequestDTO struct {
	Name *string `json:"name"`
	Row  *int    `json:"row"`
	Col  *int    `json:"col"`
}

// SheetDTO 工作表通用 DTO（用于列表和详情展示）
type SheetDTO struct {
	ID        int64  `json:"id"`         // 工作表 ID
	Name      string `json:"name"`       // 工作表名称
	CreatorID int64  `json:"creator_id"` // 创建者 ID
	ClassID   int64  `json:"class_id"`   // 班级 ID
}

// SheetListResponseDTO 工作表列表分页响应 DTO
type SheetListResponseDTO struct {
	Total  int64      `json:"total"`  // 数据总数
	Page   int        `json:"page"`   // 当前页码
	Sheets []SheetDTO `json:"sheets"` // 工作表数据列表
}

// SheetDetailResponseDTO 工作表详情 DTO，包含单元格数据
type SheetDetailResponseDTO struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Week    int    `json:"week"`
	Row     int    `json:"row"`
	Col     int    `json:"col"`
	ClassID int64  `json:"class_id"`
}
