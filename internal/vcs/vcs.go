package vcs

import (
	"fmt"
	"os/exec"
)

type VCS interface {
	Diff() (string, error)
	Commit(message string) error
}

// todo use go-git
type git struct{}

func (g *git) Diff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running git diff: %w", err)
	}

	return string(output), nil
}

func (g *git) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-am", message)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running git commit: %w", err)
	}

	return nil
}

func NewGit() VCS {
	return &git{}
}
