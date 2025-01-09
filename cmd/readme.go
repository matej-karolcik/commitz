package cmd

import (
	"commitz/internal/generate"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms/ollama"
)

var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate a README file",
	Long:  ``,
	RunE:  runReadme,
}

func runReadme(_ *cobra.Command, _ []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	files, err := os.ReadDir(wd)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	filemap := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(wd, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		filemap[file.Name()] = string(content)
	}

	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return fmt.Errorf("failed to create llm: %w", err)
	}

	readme, err := generate.Readme(context.Background(), llm, filemap)
	if err != nil {
		return fmt.Errorf("failed to generate readme: %w", err)
	}

	readmeFile, err := os.Create(filepath.Join(wd, "readme.md"))
	if err != nil {
		return fmt.Errorf("failed to create readme file: %w", err)
	}

	if _, err = readmeFile.WriteString(readme); err != nil {
		return fmt.Errorf("failed to write readme file: %w", err)
	}

	return nil
}
