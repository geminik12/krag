package v1

// ChatRequest 对话请求
type ChatRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
	Model          string `json:"model"`
	Stream         bool   `json:"stream"`
}

// ChatResponse 对话响应（非流式）
type ChatResponse struct {
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
}

// ChatStreamResponse 对话流式响应
type ChatStreamResponse struct {
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
	Done      bool   `json:"done"`
}
