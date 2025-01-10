package cmd

import (
	"context"
	"fmt"
	"github.com/matej-karolcik/commitz/internal/ai"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
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

	excludes, err := collectExcludes(wd)

	filemap, err := collectFiles(wd, excludes)
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return fmt.Errorf("failed to create llm: %w", err)
	}

	readme, err := ai.NewOllama(llm).ReadmeFile(context.Background(), filemap)
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

func collectFiles(
	workingDirectory string,
	excludes []glob.Glob,
) (map[string]string, error) {
	filemap := make(map[string]string)

	if err := filepath.WalkDir(workingDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.Contains(path, ".idea") {
			return nil
		}

		if strings.Contains(path, ".git/") {
			return nil
		}

		for _, exclude := range excludes {
			if exclude.Match(strings.TrimLeft(d.Name(), workingDirectory+"/")) {
				fmt.Printf("excluded: %s\n", d.Name())
				return nil
			}
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		filemap[path] = string(content)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return filemap, nil
}

func collectExcludes(workingDirectory string) ([]glob.Glob, error) {
	gitignore, err := os.ReadFile(filepath.Join(workingDirectory, ".gitignore"))
	if err != nil {
		return nil, nil
	}

	var excludes []glob.Glob
	lines := strings.Split(string(gitignore), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		g, err := glob.Compile(line)
		if err != nil {
			continue
		}

		excludes = append(excludes, g)
	}

	return excludes, nil
}
