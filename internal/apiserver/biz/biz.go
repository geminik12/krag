package biz

import (
	userv1 "github.com/geminik12/krag/internal/apiserver/biz/v1/user"
	"github.com/geminik12/krag/internal/apiserver/store"
)

// IBiz 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口.
	UserV1() userv1.UserBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store store.IStore
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(store store.IStore) *biz {
	return &biz{store: store}
}

// UserV1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}
