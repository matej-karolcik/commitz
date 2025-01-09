package generate

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

const promptCommitMessage = "Write a single-line commit message given the `git diff` output."

func Commit(
	ctx context.Context,
	llm llms.Model,
	diff string,
) (string, error) {
	response, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, promptCommitMessage),
			llms.TextParts(llms.ChatMessageTypeHuman, diff),
		},
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	return response.Choices[0].Content, nil
}
