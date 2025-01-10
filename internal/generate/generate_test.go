package generate_test

import (
	"context"
	"github.com/matej-karolcik/commitz/internal/generate"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms/ollama"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	diff, err := os.ReadFile("bigdiff.txt")
	assert.NoError(t, err)

	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	assert.NoError(t, err)

	commitMsg, err := generate.Commit(
		context.Background(),
		llm,
		string(diff),
	)
	assert.NoError(t, err)

	t.Log(commitMsg)
}
