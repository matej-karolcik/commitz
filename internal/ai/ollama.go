package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

const (
	commitPrompt = `Given the git changes below, please draft a concise commit message that accurately summarizes the modifications. Follow these guidelines:

	1. Limit your commit message to 10 words.
	2. The whole commit message should in lowercase, no uppercase characters are allowed.
	3. Do not respond in Markdown

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
		llms.WithMaxLength(20),
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
