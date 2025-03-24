package service

// CreateDocumentService 创建文档服务
// func CreateDocumentService(ctx context.Context, dto *DTO.CreateDocumentRequestDTO) (*DTO.CreateDocumentResponseDTO, *apiError.ApiError) {
// 	// 假设从 JWT 或上下文中获取当前用户ID，并赋值给 dto.OwnerID
// 	// dto.OwnerID = 当前用户ID
// 	doc, err := dao.CreateDocument(ctx, dto)
// 	if err != nil {
// 		return nil, &apiError.ApiError{
// 			Code: code.ServerError,
// 			Msg:  "创建文档失败",
// 		}
// 	}
// 	return &DTO.CreateDocumentResponseDTO{
// 		DocumentID: doc.ID,
// 		Title:      doc.Title,
// 		OwnerID:    doc.OwnerID,
// 		CreateTime: doc.CreateTime,
// 		UpdateTime: doc.UpdateTime,
// 	}, nil
// }
