package v1

import "mime/multipart"

// UploadKnowledgeRequest 上传知识库文件请求
type UploadKnowledgeRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// UploadKnowledgeResponse 上传响应
type UploadKnowledgeResponse struct {
	DocumentID string `json:"document_id"`
	Chunks     int    `json:"chunks"`
}

// SearchKnowledgeRequest 搜索请求（测试用）
type SearchKnowledgeRequest struct {
	Query string `json:"query" form:"query" binding:"required"`
	K     int    `json:"k" form:"k"` // Top K
}

// SearchKnowledgeResponse 搜索响应
type SearchKnowledgeResponse struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Content string  `json:"content"`
	Score   float32 `json:"score"`
}
