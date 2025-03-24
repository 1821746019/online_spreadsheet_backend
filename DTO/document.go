package DTO

import "time"

// CreateDocumentRequestDTO 创建文档请求参数
type CreateDocumentRequestDTO struct {
	Title   string `json:"title" binding:"required"`
	OwnerID int64  `json:"-"`
}

// CreateDocumentResponseDTO 创建文档响应参数
type CreateDocumentResponseDTO struct {
	DocumentID int64     `json:"document_id"`
	Title      string    `json:"title"`
	OwnerID    int64     `json:"owner_id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// DocumentDTO 文档信息
type DocumentDTO struct {
	DocumentID int64     `json:"document_id"`
	Title      string    `json:"title"`
	OwnerID    int64     `json:"owner_id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// UpdateDocumentRequestDTO 更新文档请求参数
type UpdateDocumentRequestDTO struct {
	Title string `json:"title" binding:"required"`
}

// ShareDocumentRequestDTO 文档共享请求参数
type ShareDocumentRequestDTO struct {
	UserID     int64  `json:"user_id" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}
