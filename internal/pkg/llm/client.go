package llm

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"
)

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type Message struct {
	Role    string
	Content string
}

type Client interface {
	Chat(ctx context.Context, messages []Message, model string) (string, error)
	ChatStream(ctx context.Context, messages []Message, model string, onDelta func(string) error) error
	GetEmbedder(model string) (embeddings.Embedder, error)
}
