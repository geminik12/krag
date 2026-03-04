package llm

import (
	"sync"

	"github.com/pkoukk/tiktoken-go"
)

var (
	tkm     *tiktoken.Tiktoken
	tkmOnce sync.Once
)

func getTokenizer() *tiktoken.Tiktoken {
	tkmOnce.Do(func() {
		// Use cl100k_base as a general purpose tokenizer (GPT-4)
		// It's not exact for Llama but good enough for estimation
		var err error
		tkm, err = tiktoken.GetEncoding("cl100k_base")
		if err != nil {
			// Fallback or log error?
			// For now, we assume it works or we handle nil later
			tkm = nil
		}
	})
	return tkm
}

// CountTokens returns the estimated number of tokens in the text.
func CountTokens(text string) int {
	t := getTokenizer()
	if t == nil {
		// Fallback: rough estimate
		// 1 token ~= 4 chars in English
		// 1 token ~= 1 char in Chinese (roughly)
		// We use a conservative estimate: 1 char = 1 token
		return len([]rune(text))
	}
	// encoding is thread-safe
	return len(t.Encode(text, nil, nil))
}
