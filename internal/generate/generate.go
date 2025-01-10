package generate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

const prompt = `You are a Git commit message generator. Given the following Git diff, create a meaningful commit message that summarizes the changes made.

- Provide a short summary of the changes in the first line (ideally under 50 characters).
- Ideally there is only the first line.
- Do not explain the 'why' behind these changes.
- Do not show any code.
- Use bullet points for multiple changes if necessary.
- Tone: Keep it professional and clear.

Here is the git diff:
`

const readmePrompt = `You are a README generator. Given the following files and their contents, create a README file that summarizes the project.`

func Commit(
	ctx context.Context,
	llm llms.Model,
	diff string,
) (string, error) {
	response, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
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

func Readme(
	ctx context.Context,
	llm llms.Model,
	filemap map[string]string,
) (string, error) {
	var builder strings.Builder

	for file, content := range filemap {
		builder.WriteString(fmt.Sprintf("\n\n// File: %s\n%s", file, content))
	}

	fmt.Printf("len: %d\n", len(builder.String()))

	response, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, readmePrompt),
			llms.TextParts(llms.ChatMessageTypeSystem, builder.String()),
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
