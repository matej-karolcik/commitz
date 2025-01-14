package ai_test

import (
	"context"
	"github.com/matej-karolcik/commitz/internal/ai"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms/ollama"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	diff, err := os.ReadFile("bigdiff.txt")
	assert.NoError(t, err)

	llm, err := ollama.New(ollama.WithModel("llama3.2"), ollama.WithRunnerNumCtx(8192))
	assert.NoError(t, err)

	commitMsg, err := ai.
		NewOllama(llm, 0.2).
		CommitMessage(context.Background(), string(diff))
	assert.NoError(t, err)

	t.Log(commitMsg)
}
