package v1

import "time"

// Conversation 表示会话信息
type Conversation struct {
	ID             int64     `json:"id,omitempty"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	ModelName      string    `json:"model_name"`
	SystemPrompt   string    `json:"system_prompt,omitempty"`
	Status         int       `json:"status"`
	TotalTokens    int       `json:"total_tokens"`
	MessageCount   int       `json:"message_count"`
	LastMessageAt  time.Time `json:"last_message_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Message 表示消息信息
type Message struct {
	ID             int64     `json:"id,omitempty"`
	MessageID      string    `json:"message_id"`
	ConversationID string    `json:"conversation_id"`
	ParentID       string    `json:"parent_id,omitempty"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
	Metadata       string    `json:"metadata,omitempty"`
	Tokens         int       `json:"tokens"`
	Status         int       `json:"status"`
	Sequence       int       `json:"sequence"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListConversationRequest 获取会话列表请求
type ListConversationRequest struct {
	Limit  int `form:"limit,default=10" json:"limit"`
	Offset int `form:"offset,default=0" json:"offset"`
}

// ListConversationResponse 获取会话列表响应
type ListConversationResponse struct {
	TotalCount    int64           `json:"totalCount"`
	Conversations []*Conversation `json:"conversations"`
}

// GetConversationRequest 获取单个会话详情请求
type GetConversationRequest struct {
	ConversationID string `uri:"id" binding:"required"`
}

// ListMessageRequest 获取消息列表请求
type ListMessageRequest struct {
	ConversationID string `uri:"id" binding:"required"`
	Limit          int    `form:"limit,default=20" json:"limit"`
	BeforeID       string `form:"before_id" json:"before_id"` // 分页游标，获取该message_id之前的消息
}

// ListMessageResponse 获取消息列表响应
type ListMessageResponse struct {
	Messages []*Message `json:"messages"`
	HasMore  bool       `json:"has_more"`
}

// CreateConversationRequest 创建会话请求
type CreateConversationRequest struct {
	ModelName    string `json:"model_name" binding:"required"`
	SystemPrompt string `json:"system_prompt"`
}

// CreateConversationResponse 创建会话响应
type CreateConversationResponse struct {
	ConversationID string `json:"conversation_id"`
}

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Role           string `json:"role" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

// CreateMessageResponse 创建消息响应
type CreateMessageResponse struct {
	MessageID string `json:"message_id"`
}
