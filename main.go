package main

import (
	"commitz/internal/generate"
	"commitz/internal/vcs"
	"context"
	"fmt"
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

	fmt.Printf("git commit -am\"%s\"\n", commitMsg)

	return nil
}
