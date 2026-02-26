package handler

import (
	"context"
	"encoding/json"

	"github.com/geminik12/krag/internal/pkg/contextx"
	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/geminik12/krag/internal/pkg/llm"
	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

// Chat 对话接口
func (h *Handler) Chat(c *gin.Context) {
	var r v1.ChatRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errorsx.ErrBind, nil)
		return
	}

	// 0. 确保会话存在
	userID := contextx.UserID(c.Request.Context())
	if err := h.biz.ConversationV1().CreateIfNotExist(c.Request.Context(), r.ConversationID, userID, r.Model); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	// 1. 获取历史消息
	history, err := h.biz.ConversationV1().ListMessages(c.Request.Context(), &v1.ListMessageRequest{
		ConversationID: r.ConversationID,
		Limit:          50, // 获取更多历史消息以便计算 Token
	})
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	// 2. 构建消息上下文 (Token Management)
	// 当前用户问题
	currentUserMsg := llm.Message{
		Role:    llm.RoleUser,
		Content: r.Content,
	}
	currentTokens := llm.CountTokens(currentUserMsg.Content)

	const MaxContextTokens = 3000 // 预留 1000 token 给回复和 System Prompt

	var historyMessages []llm.Message
	totalTokens := currentTokens

	// history.Messages 是按时间倒序排列的 (最新的在前)
	for _, msg := range history.Messages {
		msgTokens := llm.CountTokens(msg.Content)
		if totalTokens+msgTokens > MaxContextTokens {
			break
		}
		totalTokens += msgTokens
		historyMessages = append(historyMessages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 反转历史消息顺序，使其按时间正序排列 [oldest -> newest]
	for i, j := 0, len(historyMessages)-1; i < j; i, j = i+1, j-1 {
		historyMessages[i], historyMessages[j] = historyMessages[j], historyMessages[i]
	}

	// 组合最终的消息列表
	messages := append(historyMessages, currentUserMsg)

	// 3. 保存用户问题到数据库
	_, err = h.biz.ConversationV1().CreateMessage(context.Background(), &v1.CreateMessageRequest{
		ConversationID: r.ConversationID,
		Role:           llm.RoleUser,
		Content:        r.Content,
	})
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	// 4. 调用 LLM
	if r.Stream {
		h.handleStreamChat(c, &r, messages)
	} else {
		h.handleNormalChat(c, &r, messages)
	}
}

func (h *Handler) handleNormalChat(c *gin.Context, r *v1.ChatRequest, messages []llm.Message) {
	respContent, err := h.biz.ChatV1().Chat(c.Request.Context(), messages, r.Model)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	// 保存 AI 回复到数据库
	aiMsg, err := h.biz.ConversationV1().CreateMessage(context.Background(), &v1.CreateMessageRequest{
		ConversationID: r.ConversationID,
		Role:           llm.RoleAssistant,
		Content:        respContent,
	})

	var messageID string
	if err == nil {
		messageID = aiMsg.MessageID
	}

	core.WriteResponse(c, &v1.ChatResponse{
		MessageID: messageID,
		Content:   respContent,
	}, nil)
}

func (h *Handler) handleStreamChat(c *gin.Context, r *v1.ChatRequest, messages []llm.Message) {
	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	var fullContent string

	err := h.biz.ChatV1().ChatStream(c.Request.Context(), messages, r.Model, func(chunk string) error {
		fullContent += chunk
		resp := &v1.ChatStreamResponse{
			Content: chunk,
			Done:    false,
		}
		data, _ := json.Marshal(resp)
		c.SSEvent("message", string(data))
		c.Writer.Flush()
		return nil
	})

	if err != nil {
		// 流式传输中出错，尝试发送错误事件
		c.SSEvent("error", err.Error())
		c.Writer.Flush()
	}

	// 发送结束标志
	resp := &v1.ChatStreamResponse{
		Done: true,
	}
	data, _ := json.Marshal(resp)
	c.SSEvent("message", string(data))
	c.Writer.Flush()

	// 保存 AI 完整回复到数据库
	if len(fullContent) > 0 {
		h.biz.ConversationV1().CreateMessage(context.Background(), &v1.CreateMessageRequest{
			ConversationID: r.ConversationID,
			Role:           llm.RoleAssistant,
			Content:        fullContent,
		})
	}
}
