package dao

// func CreateDocument(ctx context.Context, dto *DTO.CreateDocumentRequestDTO) (*model.Document, error) {
// 	doc := model.Document{
// 		Title:      dto.Title,
// 		OwnerID:    dto.OwnerID,
// 		CreateTime: time.Now(),
// 		UpdateTime: time.Now(),
// 		DeleteTime: 0, // 0 表示未删除
// 	}
// 	if err := mysql.GetDB().WithContext(ctx).Create(&doc).Error; err != nil {
// 		return nil, err
// 	}
// 	return &doc, nil
// }
