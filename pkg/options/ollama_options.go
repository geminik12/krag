package options

import "github.com/geminik12/krag/internal/pkg/llm"

type OllamaOptions struct {
	Host         string `json:"host,omitempty" mapstructure:"host"`
	DefaultModel string `json:"default-model,omitempty" mapstructure:"default-model"`
	SystemPrompt string `json:"system-prompt,omitempty" mapstructure:"system-prompt"`
}

func NewOllamaOptions() *OllamaOptions {
	return &OllamaOptions{
		Host:         "http://localhost:11434",
		DefaultModel: "llama3.2",
		SystemPrompt: "You are a helpful assistant.",
	}
}

func (o *OllamaOptions) Validate() error {
	return nil
}

func (o *OllamaOptions) NewClient() llm.Client {
	return llm.NewOllamaClient(
		llm.WithHost(o.Host),
		llm.WithDefaultModel(o.DefaultModel),
		llm.WithSystemPrompt(o.SystemPrompt),
	)
}
