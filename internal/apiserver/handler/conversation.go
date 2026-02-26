package handler

import (
	"github.com/geminik12/krag/internal/pkg/contextx"
	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

// CreateConversation 创建会话
func (h *Handler) CreateConversation(c *gin.Context) {
	var r v1.CreateConversationRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errorsx.ErrBind, nil)
		return
	}

	username := contextx.UserID(c.Request.Context())
	resp, err := h.biz.ConversationV1().Create(c.Request.Context(), username, &r)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// ListConversation 获取会话列表
func (h *Handler) ListConversation(c *gin.Context) {
	var r v1.ListConversationRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	username := contextx.UserID(c.Request.Context())
	resp, err := h.biz.ConversationV1().List(c.Request.Context(), username, &r)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// GetConversation 获取单个会话详情
func (h *Handler) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	resp, err := h.biz.ConversationV1().Get(c.Request.Context(), conversationID)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// DeleteConversation 删除会话
func (h *Handler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	if err := h.biz.ConversationV1().Delete(c.Request.Context(), conversationID); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, nil, nil)
}

// ListMessage 获取消息列表
func (h *Handler) ListMessage(c *gin.Context) {
	var r v1.ListMessageRequest
	r.ConversationID = c.Param("id")
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	resp, err := h.biz.ConversationV1().ListMessages(c.Request.Context(), &r)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}
