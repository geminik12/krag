package chat

import (
	"context"

	"github.com/geminik12/krag/internal/pkg/llm"
)

type ChatBiz interface {
	Chat(ctx context.Context, messages []llm.Message, model string) (string, error)
	ChatStream(ctx context.Context, messages []llm.Message, model string, onDelta func(string) error) error
}

type chatBiz struct {
	client llm.Client
}

func New(client llm.Client) ChatBiz {
	return &chatBiz{client: client}
}

func (b *chatBiz) Chat(ctx context.Context, messages []llm.Message, model string) (string, error) {
	messages = b.ensureSystemPrompt(messages)
	return b.client.Chat(ctx, messages, model)
}

func (b *chatBiz) ChatStream(ctx context.Context, messages []llm.Message, model string, onDelta func(string) error) error {
	messages = b.ensureSystemPrompt(messages)
	return b.client.ChatStream(ctx, messages, model, onDelta)
}

func (b *chatBiz) ensureSystemPrompt(messages []llm.Message) []llm.Message {
	for _, msg := range messages {
		if msg.Role == llm.RoleSystem {
			return messages
		}
	}

	systemMsg := llm.Message{
		Role:    llm.RoleSystem,
		Content: "You are a helpful assistant. Please always answer in Chinese.",
	}
	return append([]llm.Message{systemMsg}, messages...)
}
