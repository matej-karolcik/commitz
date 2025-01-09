package generate

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

const promptCommitMessage = `You are an assistant to write a commit message. The user will send you the content of the commit diff, and you will reply with a concise, single-line commit message without any explanation.`

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
