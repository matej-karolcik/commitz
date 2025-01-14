package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/matej-karolcik/commitz/internal/ai"
	"github.com/matej-karolcik/commitz/internal/config"
	"github.com/matej-karolcik/commitz/internal/vcs"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms/ollama"
)

var (
	prefix     string
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "Prefix for the commit message")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to the config file")
}

var rootCmd = &cobra.Command{
	Use:   "commitz",
	Short: "A brief description of your application",
	Long:  ``,
	RunE:  run,
}

func run(*cobra.Command, []string) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("creating ollama client: %w", err)
	}

	_, err = ollamaClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing models: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel(cfg.Model), ollama.WithRunnerNumCtx(8192))
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

	commitMsg, err := ai.NewOllama(llm, cfg.Temperature).CommitMessage(
		ctx,
		diff,
	)
	if err != nil {
		return fmt.Errorf("generating commit message: %w", err)
	}

	if prefix != "" {
		commitMsg = prefix + " " + commitMsg
	}

	fmt.Println(commitMsg)

	prompt := promptui.Prompt{
		Label:     "Commit with this message?",
		Default:   "y",
		IsConfirm: true,
	}

	if _, err = prompt.Run(); err != nil {
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
