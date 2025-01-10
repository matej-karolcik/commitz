package cmd

import (
	"commitz/internal/generate"
	"commitz/internal/vcs"
	"context"
	"errors"
	"fmt"
	"github.com/ollama/ollama/api"
	"github.com/tmc/langchaingo/llms/ollama"
	"os"
	"os/signal"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(readmeCmd)
}

var rootCmd = &cobra.Command{
	Use:   "commitz",
	Short: "A brief description of your application",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			if args[0] == "readme" {
				return readmeCmd.RunE(cmd, args[1:])
			}
		}
		return run(cmd, args)
	},
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

	if prefix != "" {
		commitMsg = prefix + " " + commitMsg
	}

	fmt.Println(commitMsg)

	prompt := promptui.Prompt{
		Label:     "Commit with this message?",
		Default:   commitMsg,
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
