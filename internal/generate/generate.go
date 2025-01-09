package generate

import (
	"context"
	"errors"
	"fmt"
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

func Commit(
	ctx context.Context,
	llm llms.Model,
	diff string,
) (string, error) {
	response, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt+diff),
			//llms.TextParts(llms.ChatMessageTypeHuman, diff),
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
