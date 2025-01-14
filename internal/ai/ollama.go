package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

const (
	commitPrompt = `Given the git changes below, please draft a concise commit message that accurately summarizes the modifications. Follow these guidelines:

1. Be concise, with a strict limit of 10 words.
2. Begin with an appropriate prefix indicating the type of change (e.g., feat: for new features, fix: for bug fixes, refactor: for code refactoring, chore: for routine tasks).
3. You go deeper in the changes made and understand what changed.
4. Never respond in markdown format.
5. You give meaningful commit message.
6. Never respond with anything else than the commit message.
7. You must respond with a single line commit message.
8. You must respond in lowercase.
9. You must response in present tense.
10. You must not provide any comments or explanations.

Git Changes: 

`
)

var _ AI = (*ollama)(nil)

type ollama struct {
	backend     llms.Model
	temperature float64
}

func NewOllama(backend llms.Model, temperature float64) AI {
	return &ollama{
		backend:     backend,
		temperature: temperature,
	}
}

func (o *ollama) CommitMessage(
	ctx context.Context,
	diff string,
) (string, error) {
	response, err := o.ask(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
			llms.TextParts(llms.ChatMessageTypeHuman, commitPrompt+diff),
		},
		llms.WithTemperature(o.temperature),
		llms.WithMaxLength(10),
		// llms.WithTopK(1),
		llms.WithTopP(0.1),
		llms.WithFrequencyPenalty(2.0),
		// llms.WithPresencePenalty(2.0),
	)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
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
