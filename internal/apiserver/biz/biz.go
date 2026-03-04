package biz

import (
	chatv1 "github.com/geminik12/krag/internal/apiserver/biz/v1/chat"
	conversationv1 "github.com/geminik12/krag/internal/apiserver/biz/v1/conversation"
	knowledgev1 "github.com/geminik12/krag/internal/apiserver/biz/v1/knowledge"
	userv1 "github.com/geminik12/krag/internal/apiserver/biz/v1/user"
	"github.com/geminik12/krag/internal/apiserver/store"
	"github.com/geminik12/krag/internal/pkg/llm"
	"github.com/geminik12/krag/internal/pkg/storage"
	genericoptions "github.com/geminik12/krag/pkg/options"
	"github.com/hibiken/asynq"
)

// IBiz 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口.
	UserV1() userv1.UserBiz
	ConversationV1() conversationv1.ConversationBiz
	ChatV1() chatv1.ChatBiz
	KnowledgeV1() knowledgev1.KnowledgeBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store         store.IStore
	client        llm.Client
	storage       storage.Storage
	qdrantOptions *genericoptions.QdrantOptions
	asynqClient   *asynq.Client
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(store store.IStore, client llm.Client, storage storage.Storage, qdrantOptions *genericoptions.QdrantOptions, asynqClient *asynq.Client) *biz {
	return &biz{
		store:         store,
		client:        client,
		storage:       storage,
		qdrantOptions: qdrantOptions,
		asynqClient:   asynqClient,
	}
}

// UserV1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}

func (b *biz) ConversationV1() conversationv1.ConversationBiz {
	return conversationv1.New(b.store)
}

func (b *biz) ChatV1() chatv1.ChatBiz {
	return chatv1.New(b.client)
}

func (b *biz) KnowledgeV1() knowledgev1.KnowledgeBiz {
	return knowledgev1.New(b.client, b.storage, b.store, b.qdrantOptions, b.asynqClient)
}
