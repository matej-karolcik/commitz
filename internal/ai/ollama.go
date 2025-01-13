package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

const (
	commitPrompt = `You are a Git commit message generator. Given the following Git diff, create a meaningful commit message that summarizes the changes made.

- Provide a short summary of the most significant changes (ideally under 20 words).
- Use present tense.
- You must not include any feedback, suggestions, or code snippets.
- Ideally there is only the first line.
- Do not explain the 'why' behind these changes.
- Do not show any code.
- Tone: Keep it professional and clear.

Here is the git diff:
`

	summarizeFilePrompt  = "Summarize this file, the summary should be used to generate a readme in another prompt. Do not provide any comments or ask any questions. The summary should be short and it should state what the file does."
	generateReadmePrompt = "Can you generate a readme.md from this list of file summaries?"
)

var _ AI = (*ollama)(nil)

type ollama struct {
	backend llms.Model
}

func NewOllama(backend llms.Model) AI {
	return &ollama{
		backend: backend,
	}
}

func (o *ollama) CommitMessage(
	ctx context.Context,
	diff string,
) (string, error) {
	response, err := o.ask(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, commitPrompt),
			llms.TextParts(llms.ChatMessageTypeHuman, diff),
		},
		llms.WithTemperature(0.5),
		llms.WithMaxLength(20),
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	return response, nil
}

func (o *ollama) ReadmeFile(
	ctx context.Context,
	filemap map[string]string,
) (string, error) {
	var builder strings.Builder

	summaries := make(map[string]string, len(filemap))
	for file, content := range filemap {
		summary, err := o.summarizeFile(ctx, file, content)
		if err != nil {
			return "", fmt.Errorf("summarizing file: %w", err)
		}

		summaries[file] = summary
	}

	for file, summary := range summaries {
		builder.WriteString(fmt.Sprintf("\n\n// File: %s\n%s", file, summary))
	}

	response, err := o.ask(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, generateReadmePrompt),
			llms.TextParts(llms.ChatMessageTypeHuman, builder.String()),
		},
	)
	if err != nil {
		return "", fmt.Errorf("asking llm: %w", err)
	}

	return response, nil
}

func (o *ollama) summarizeFile(
	ctx context.Context,
	filename, content string,
) (string, error) {
	response, err := o.ask(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, summarizeFilePrompt),
			llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf("Filename: %s\nContent: %s\n\n", filename, content)),
		},
	)
	if err != nil {
		return "", fmt.Errorf("asking llm: %w", err)
	}

	return response, nil
}

func (o *ollama) ask(
	ctx context.Context,
	prompts []llms.MessageContent,
	callOptions ...llms.CallOption,
) (string, error) {
	response, err := o.backend.GenerateContent(ctx, prompts, callOptions...)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response were provided by llm")
	}

	return response.Choices[0].Content, nil
}
