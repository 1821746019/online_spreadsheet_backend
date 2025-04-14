package DTO

type CellDTO struct {
	ID        int64  `json:"id"`
	SheetID   int64  `json:"sheet_id"`
	RowIndex  int    `json:"row_index"`
	ColIndex  int    `json:"col_index"`
	ItemID    *int64 `json:"item_id"`
	Content   string `json:"content"`
	ClassName string `json:"class_name"`
	WeekType  string `json:"week_type"`
	ClassRoom string `json:"class_room"`
}

type DeleteItemInCellRequest struct {
	Row int `json:"row", binding:"required"`
	Col int `json:"col", binding:"required"`
}

type CreateDragItemRequestDTO struct {
	Content          string  `json:"content" binding:"required"`
	WeekType         string  `json:"week_type" binding:"required"`
	Classroom        string  `json:"classroom" binding:"required"`
	SelectedClassIDs []int64 `json:"selected_class_ids,required"` // 使用omitempty实现可选
}

type UpdateDragItemRequestDTO struct {
	Content          string  `json:"content"`
	WeekType         string  `json:"week_type"`
	Classroom        string  `json:"classroom"`
	SelectedClassIDs []int64 `json:"selected_class_ids"`
}

type DragItemResponseDTO struct {
	ID         int64    `json:"id"`
	WeekType   string   `json:"week_type"`
	Classroom  string   `json:"classroom"`
	ClassNames []string `json:"class_names"`
	Teacher    string   `json:"teacher"`
	Content    string   `json:"content"`
	CreatorID  int64    `json:"creator_id"`
	CreateTime string   `json:"create_time"`
	UpdateTime string   `json:"update_time"`
}

type MoveDragItemRequest struct {
	TargetRow int `json:"target_row" binding:"required"`
	TargetCol int `json:"target_col" binding:"required"`
}
