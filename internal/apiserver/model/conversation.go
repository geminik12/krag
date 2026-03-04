package model

import "time"

const (
	TableNameConversation = "conversations"
	TableNameMessage      = "messages"
)

// Conversation 会话模型
type Conversation struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ConversationID string    `gorm:"column:conversation_id;not null;unique;comment:会话唯一ID" json:"conversation_id"`
	UserID         string    `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`
	Title          string    `gorm:"column:title;not null;default:'';comment:会话标题" json:"title"`
	ModelName      string    `gorm:"column:model_name;not null;default:'';comment:使用的模型名称" json:"model_name"`
	SystemPrompt   string    `gorm:"column:system_prompt;type:text;comment:系统提示词" json:"system_prompt"`
	Status         int       `gorm:"column:status;not null;default:1;comment:状态：1-活跃，2-归档，3-删除" json:"status"`
	TotalTokens    int       `gorm:"column:total_tokens;not null;default:0;comment:累计消耗token数" json:"total_tokens"`
	MessageCount   int       `gorm:"column:message_count;not null;default:0;comment:消息总数" json:"message_count"`
	LastMessageAt  time.Time `gorm:"column:last_message_at;not null;default:current_timestamp;comment:最后消息时间" json:"last_message_at"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:current_timestamp" json:"updated_at"`
}

// TableName 表名
func (*Conversation) TableName() string {
	return TableNameConversation
}

// Message 消息模型
type Message struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MessageID      string    `gorm:"column:message_id;not null;unique;comment:消息唯一ID" json:"message_id"`
	ConversationID string    `gorm:"column:conversation_id;not null;comment:所属会话ID" json:"conversation_id"`
	ParentID       string    `gorm:"column:parent_id;default:null;comment:父消息ID" json:"parent_id"`
	Role           string    `gorm:"column:role;not null;comment:角色" json:"role"`
	Content        string    `gorm:"column:content;not null;type:longtext;comment:消息内容" json:"content"`
	ContentType    string    `gorm:"column:content_type;not null;default:'text';comment:内容类型" json:"content_type"`
	Metadata       *string   `gorm:"column:metadata;type:json;comment:元数据" json:"metadata"` // GORM 会自动处理 JSON，但这里先用 string 存，或者可以用 datatypes.JSON
	Tokens         int       `gorm:"column:tokens;not null;default:0;comment:本条消息消耗token数" json:"tokens"`
	Status         int       `gorm:"column:status;not null;default:1;comment:状态：1-正常，2-撤回，3-已编辑" json:"status"`
	Sequence       int       `gorm:"column:sequence;not null;default:0;comment:会话内序号" json:"sequence"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:current_timestamp" json:"updated_at"`
}

// TableName 表名
func (*Message) TableName() string {
	return TableNameMessage
}
