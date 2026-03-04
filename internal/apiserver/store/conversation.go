package store

import (
	"context"

	"github.com/geminik12/krag/internal/apiserver/model"
	"gorm.io/gorm"
)

type ConversationStore interface {
	Create(ctx context.Context, conversation *model.Conversation) error
	Get(ctx context.Context, conversationID string) (*model.Conversation, error)
	List(ctx context.Context, userID string, limit, offset int) ([]*model.Conversation, int64, error)
	Delete(ctx context.Context, conversationID string) error

	// Message operations
	CreateMessage(ctx context.Context, message *model.Message) error
	ListMessages(ctx context.Context, conversationID string, limit int, beforeID string) ([]*model.Message, error)
}

type conversationStore struct {
	db *gorm.DB
}

func newConversations(db *gorm.DB) *conversationStore {
	return &conversationStore{db: db}
}

func (s *conversationStore) Create(ctx context.Context, conversation *model.Conversation) error {
	return s.db.Create(conversation).Error
}

func (s *conversationStore) Get(ctx context.Context, conversationID string) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := s.db.Where("conversation_id = ?", conversationID).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (s *conversationStore) List(ctx context.Context, userID string, limit, offset int) ([]*model.Conversation, int64, error) {
	var conversations []*model.Conversation
	var count int64

	db := s.db.Model(&model.Conversation{}).Where("user_id = ? AND status = 1", userID)

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("last_message_at DESC").Limit(limit).Offset(offset).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, count, nil
}

func (s *conversationStore) Delete(ctx context.Context, conversationID string) error {
	// 软删除，状态改为 3
	return s.db.Model(&model.Conversation{}).Where("conversation_id = ?", conversationID).Update("status", 3).Error
}

func (s *conversationStore) CreateMessage(ctx context.Context, message *model.Message) error {
	return s.db.Create(message).Error
}

func (s *conversationStore) ListMessages(ctx context.Context, conversationID string, limit int, beforeID string) ([]*model.Message, error) {
	var messages []*model.Message
	// 这里必须是 s.db.Model(&model.Message{}) 否则 GORM 不知道表名
	db := s.db.Model(&model.Message{}).Where("conversation_id = ? AND status = 1", conversationID)

	if beforeID != "" {
		// 假设 sequence 是递增的，可以用 sequence 做游标，或者直接用 id/created_at
		// 这里简化处理，如果传入了 beforeID，先查出该消息的 sequence/id
		var msg model.Message
		if err := s.db.Where("message_id = ?", beforeID).First(&msg).Error; err == nil {
			db = db.Where("id < ?", msg.ID)
		}
	}

	// 按时间倒序查，展示时再正序
	if err := db.Order("id DESC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
