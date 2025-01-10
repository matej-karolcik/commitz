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
const readmePrompt2 = `You are a README generator. Given the following files and their summaries, create a README file that summarizes the project.`

func Commit(
	ctx context.Context,
	llm llms.Model,
	diff string,
) (string, error) {
	response, err := ask(
		ctx,
		llm,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
			llms.TextParts(llms.ChatMessageTypeHuman, diff),
		},
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	return response, nil
}

func Readme(
	ctx context.Context,
	llm llms.Model,
	filemap map[string]string,
) (string, error) {
	var builder strings.Builder

	summaries := make(map[string]string, len(filemap))
	for file, content := range filemap {
		summary, err := summarizeFile(ctx, llm, file, content)
		if err != nil {
			return "", fmt.Errorf("summarizing file: %w", err)
		}

		summaries[file] = summary
	}

	for file, summary := range summaries {
		builder.WriteString(fmt.Sprintf("\n\n// File: %s\n%s", file, summary))
	}

	response, err := ask(
		ctx,
		llm,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, "Can you generate a readme.md from this list of file summaries?"),
			llms.TextParts(llms.ChatMessageTypeHuman, builder.String()),
		},
	)
	if err != nil {
		return "", fmt.Errorf("asking llm: %w", err)
	}

	return response, nil
}

func summarizeFile(
	ctx context.Context,
	llm llms.Model,
	filename, content string,
) (string, error) {
	response, err := ask(
		ctx,
		llm,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "Summarize this file, the summary should be used to generate a readme in another prompt. Do not provide any comments or ask any questions. The summary should be short and it should state what the file does."),
			llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf("Filename: %s\nContent: %s\n\n", filename, content)),
		},
	)
	if err != nil {
		return "", fmt.Errorf("asking llm: %w", err)
	}

	return response, nil
}

func ask(ctx context.Context, llm llms.Model, prompts []llms.MessageContent) (string, error) {
	response, err := llm.GenerateContent(
		ctx,
		prompts,
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response were provided by llm")
	}

	return response.Choices[0].Content, nil
}
