package DTO

type CreateClassRequestDTO struct {
	Name        string `json:"name" binding:"required"` // 班级名称
	CreaetSheet bool   `json:"create_sheet"`            // 是否创建周表格
	Weeks       int    `json:"weeks"`                   // 周数
}

type UpdateClassRequestDTO struct {
	Name string `json:"name"` // 班级名称
}

type ClassResponseDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ClassListDTO struct {
	Total int64                `json:"total"`
	List  []ClassSimpleItemDTO `json:"list"`
}

type ClassSimpleItemDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
