package generate

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

const promptCommitMessage = "Write a single-line commit message given the `git diff` output. Follow best practices and do not enclose the message in quotes. Do not give any explanations."

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

	for _, choice := range response.Choices {
		fmt.Println(choice.Content)
	}

	return response.Choices[0].Content, nil
}
