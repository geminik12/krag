package handler

import (
	biz "github.com/geminik12/krag/internal/apiserver/biz"
	"github.com/geminik12/krag/internal/apiserver/pkg/validation"
)

// Handler 处理博客模块的请求.
type Handler struct {
	biz biz.IBiz
	val *validation.Validator
}

// NewHandler 创建新的 Handler 实例.
func NewHandler(biz biz.IBiz, val *validation.Validator) *Handler {
	return &Handler{
		biz: biz,
		val: val,
	}
}
