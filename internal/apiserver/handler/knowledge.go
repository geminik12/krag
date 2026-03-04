package handler

import (
	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

// UploadKnowledge 上传知识库文件
func (h *Handler) UploadKnowledge(c *gin.Context) {
	var r v1.UploadKnowledgeRequest
	// ShouldBind 会自动根据 Content-Type 选择绑定器，multipart/form-data 会绑定到 struct
	if err := c.ShouldBind(&r); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrBind.WithMessage("%v", err))
		return
	}

	resp, err := h.biz.KnowledgeV1().Upload(c.Request.Context(), &r)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// SearchKnowledge 搜索知识库
func (h *Handler) SearchKnowledge(c *gin.Context) {
	var r v1.SearchKnowledgeRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	resp, err := h.biz.KnowledgeV1().Search(c.Request.Context(), &r)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}
