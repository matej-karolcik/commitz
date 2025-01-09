package cmd

import (
	"commitz/internal/generate"
	"commitz/internal/vcs"
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/ollama/ollama/api"
	"github.com/tmc/langchaingo/llms/ollama"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commitz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var rootCmd = &cobra.Command{
	Use:   "commitz",
	Short: "A brief description of your application",
	Long:  ``,
	RunE:  run,
}

func run(_ *cobra.Command, args []string) error {
	var prefix string
	if len(args) > 0 {
		prefix = args[0]
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("creating ollama client: %w", err)
	}

	_, err = ollamaClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing models: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return fmt.Errorf("creating ollama: %w", err)
	}

	git := vcs.NewGit()

	diff, err := git.Diff()
	if err != nil {
		return fmt.Errorf("getting diff: %w", err)
	}

	if diff == "" {
		return errors.New("no changes to commit")
	}

	commitMsg, err := generate.Commit(
		ctx,
		llm,
		diff,
	)
	if err != nil {
		return fmt.Errorf("generating commit message: %w", err)
	}

	commitMsg = prefix + " " + commitMsg

	var proceed bool
	if err = huh.NewConfirm().
		Title("Press enter to commit with this message:\n\n" + commitMsg + "\n").
		Value(&proceed).
		WithTheme(huh.ThemeBase()).
		Run(); err != nil || !proceed {
		return nil
	}

	if err = git.Commit(commitMsg); err != nil {
		return fmt.Errorf("committing: %w", err)
	}

	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
