package conversation

import (
	"context"

	"github.com/geminik12/krag/internal/apiserver/model"
	"github.com/geminik12/krag/internal/apiserver/store"
	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
	"github.com/google/uuid"
)

type ConversationBiz interface {
	Create(ctx context.Context, username string, r *v1.CreateConversationRequest) (*v1.CreateConversationResponse, error)
	List(ctx context.Context, username string, r *v1.ListConversationRequest) (*v1.ListConversationResponse, error)
	Get(ctx context.Context, conversationID string) (*v1.Conversation, error)
	CreateIfNotExist(ctx context.Context, conversationID string, userID string, modelName string) error
	Delete(ctx context.Context, conversationID string) error

	// Message related
	CreateMessage(ctx context.Context, r *v1.CreateMessageRequest) (*v1.CreateMessageResponse, error)
	ListMessages(ctx context.Context, r *v1.ListMessageRequest) (*v1.ListMessageResponse, error)
}

type conversationBiz struct {
	store store.IStore
}

func New(store store.IStore) *conversationBiz {
	return &conversationBiz{store: store}
}

func (b *conversationBiz) Create(ctx context.Context, username string, r *v1.CreateConversationRequest) (*v1.CreateConversationResponse, error) {
	conversation := &model.Conversation{
		ConversationID: uuid.New().String(),
		UserID:         username, // 这里假设 username 就是 userID，实际场景可能需要转换
		ModelName:      r.ModelName,
		SystemPrompt:   r.SystemPrompt,
		Status:         1,
	}

	if err := b.store.Conversation().Create(ctx, conversation); err != nil {
		return nil, err
	}

	return &v1.CreateConversationResponse{ConversationID: conversation.ConversationID}, nil
}

func (b *conversationBiz) List(ctx context.Context, username string, r *v1.ListConversationRequest) (*v1.ListConversationResponse, error) {
	conversations, count, err := b.store.Conversation().List(ctx, username, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	var res []*v1.Conversation
	for _, c := range conversations {
		res = append(res, &v1.Conversation{
			ID:             c.ID,
			ConversationID: c.ConversationID,
			UserID:         c.UserID,
			Title:          c.Title,
			ModelName:      c.ModelName,
			SystemPrompt:   c.SystemPrompt,
			Status:         c.Status,
			TotalTokens:    c.TotalTokens,
			MessageCount:   c.MessageCount,
			LastMessageAt:  c.LastMessageAt,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		})
	}

	return &v1.ListConversationResponse{TotalCount: count, Conversations: res}, nil
}

func (b *conversationBiz) Get(ctx context.Context, conversationID string) (*v1.Conversation, error) {
	c, err := b.store.Conversation().Get(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return &v1.Conversation{
		ID:             c.ID,
		ConversationID: c.ConversationID,
		UserID:         c.UserID,
		Title:          c.Title,
		ModelName:      c.ModelName,
		SystemPrompt:   c.SystemPrompt,
		Status:         c.Status,
		TotalTokens:    c.TotalTokens,
		MessageCount:   c.MessageCount,
		LastMessageAt:  c.LastMessageAt,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}, nil
}

func (b *conversationBiz) CreateIfNotExist(ctx context.Context, conversationID string, userID string, modelName string) error {
	_, err := b.store.Conversation().Get(ctx, conversationID)
	if err == nil {
		return nil // exists
	}

	// Assume error is "not found"
	conversation := &model.Conversation{
		ConversationID: conversationID,
		UserID:         userID,
		ModelName:      modelName,
		Status:         1,
	}

	return b.store.Conversation().Create(ctx, conversation)
}

func (b *conversationBiz) Delete(ctx context.Context, conversationID string) error {
	return b.store.Conversation().Delete(ctx, conversationID)
}

func (b *conversationBiz) CreateMessage(ctx context.Context, r *v1.CreateMessageRequest) (*v1.CreateMessageResponse, error) {
	message := &model.Message{
		MessageID:      uuid.New().String(),
		ConversationID: r.ConversationID,
		Role:           r.Role,
		Content:        r.Content,
		Status:         1,
	}

	if err := b.store.Conversation().CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	return &v1.CreateMessageResponse{MessageID: message.MessageID}, nil
}

func (b *conversationBiz) ListMessages(ctx context.Context, r *v1.ListMessageRequest) (*v1.ListMessageResponse, error) {
	messages, err := b.store.Conversation().ListMessages(ctx, r.ConversationID, r.Limit, r.BeforeID)
	if err != nil {
		return nil, err
	}

	var res []*v1.Message
	for _, m := range messages {
		res = append(res, &v1.Message{
			ID:             m.ID,
			MessageID:      m.MessageID,
			ConversationID: m.ConversationID,
			Role:           m.Role,
			Content:        m.Content,
			ContentType:    m.ContentType,
			Status:         m.Status,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		})
	}

	return &v1.ListMessageResponse{Messages: res, HasMore: len(messages) == r.Limit}, nil
}
