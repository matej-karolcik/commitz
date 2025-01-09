package generate

import (
	"context"
	"errors"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

const promptCommitMessage = "Write a single-line commit message given the `git diff` output. Follow best practices and do not enclose the message in quotes. Do not give any explanations."

const prompt = `You are a Git commit message generator. Given the following Git diff, create a meaningful commit message that summarizes the changes made.

1. Provide a short summary of the changes in the first line (ideally under 50 characters).
2. List the files that were changed and describe what was changed.
3. Explain the 'why' behind these changes.
4. Use bullet points for multiple changes if applicable.
5. Tone: Keep it professional and clear.`

func Commit(
	ctx context.Context,
	llm llms.Model,
	diff string,
) (string, error) {
	response, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, prompt),
			llms.TextParts(llms.ChatMessageTypeHuman, diff),
		},
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response were provided by llm")
	}

	return response.Choices[0].Content, nil
}
