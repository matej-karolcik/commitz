package main

import (
	"commitz/internal/generate"
	"commitz/internal/vcs"
	"context"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/tmc/langchaingo/llms/ollama"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return fmt.Errorf("creating ollama: %w", err)
	}

	diff, err := vcs.NewGit().Diff()
	if err != nil {
		return fmt.Errorf("getting diff: %w", err)
	}

	if diff == "" {
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	commitMsg, err := generate.Commit(
		ctx,
		llm,
		diff,
	)
	if err != nil {
		return fmt.Errorf("generating commit message: %w", err)
	}

	var proceed bool
	if err = huh.NewConfirm().
		Title("Press enter to commit with this message:\n\n" + commitMsg + "\n").
		Value(&proceed).
		WithTheme(huh.ThemeBase()).
		Run(); err != nil || !proceed {
		return nil
	}

	if err = vcs.NewGit().Commit(commitMsg); err != nil {
		return fmt.Errorf("committing: %w", err)
	}

	return nil
}
