package llm

import (
	"context"
	"sync"

	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ollamaClient struct {
	Host         string
	DefaultModel string
	SystemPrompt string

	mu   sync.Mutex
	llms map[string]*ollama.LLM
}

type Option func(*ollamaClient)

func WithHost(host string) Option {
	return func(o *ollamaClient) {
		o.Host = host
	}
}

func WithDefaultModel(model string) Option {
	return func(o *ollamaClient) {
		o.DefaultModel = model
	}
}

func WithSystemPrompt(prompt string) Option {
	return func(o *ollamaClient) {
		o.SystemPrompt = prompt
	}
}

func NewOllamaClient(opts ...Option) Client {
	c := &ollamaClient{
		llms: make(map[string]*ollama.LLM),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (o *ollamaClient) resolveModel(model string) string {
	effectiveModel := model
	if len(effectiveModel) == 0 {
		effectiveModel = o.DefaultModel
	}

	return effectiveModel
}

func (o *ollamaClient) Chat(ctx context.Context, messages []Message, model string) (string, error) {
	effectiveModel := o.resolveModel(model)

	llm, err := o.getLLM(effectiveModel)
	if err != nil {
		return "", err
	}

	content := o.buildContent(messages)
	resp, err := llm.GenerateContent(ctx, content)
	if err != nil {
		return "", errorsx.ErrInternal.WithMessage("failed to call Ollama: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", errorsx.ErrInternal.WithMessage("empty response from Ollama")
	}

	return resp.Choices[0].Content, nil
}

func (o *ollamaClient) ChatStream(ctx context.Context, messages []Message, model string, onDelta func(string) error) error {
	effectiveModel := o.resolveModel(model)

	llm, err := o.getLLM(effectiveModel)
	if err != nil {
		return err
	}

	content := o.buildContent(messages)
	_, err = llm.GenerateContent(
		ctx,
		content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				return nil
			}
			return onDelta(string(chunk))
		}),
	)
	if err != nil {
		return errorsx.ErrInternal.WithMessage("failed to call Ollama: %v", err)
	}

	return nil
}

func (o *ollamaClient) buildContent(messages []Message) []llms.MessageContent {
	var content []llms.MessageContent

	if o.SystemPrompt != "" {
		hasSystem := false
		if len(messages) > 0 && messages[0].Role == RoleSystem {
			hasSystem = true
		}
		if !hasSystem {
			content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, o.SystemPrompt))
		}
	}

	for _, msg := range messages {
		var role llms.ChatMessageType
		switch msg.Role {
		case RoleSystem:
			role = llms.ChatMessageTypeSystem
		case RoleUser:
			role = llms.ChatMessageTypeHuman
		case RoleAssistant:
			role = llms.ChatMessageTypeAI
		default:
			role = llms.ChatMessageTypeHuman
		}
		content = append(content, llms.TextParts(role, msg.Content))
	}

	return content
}

func (o *ollamaClient) getLLM(model string) (*ollama.LLM, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(model) == 0 {
		model = o.DefaultModel
	}

	if llm, ok := o.llms[model]; ok {
		return llm, nil
	}

	opts := []ollama.Option{
		ollama.WithModel(model),
	}
	if len(o.Host) != 0 {
		opts = append(opts, ollama.WithServerURL(o.Host))
	}

	llm, err := ollama.New(opts...)
	if err != nil {
		return nil, errorsx.ErrInternal.WithMessage("failed to create Ollama client: %v", err)
	}
	o.llms[model] = llm

	return llm, nil
}
