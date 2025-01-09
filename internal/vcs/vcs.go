package vcs

import (
	"fmt"
	"os/exec"
)

type VCS interface {
	Diff() (string, error)
}

type git struct{}

func (g *git) Diff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running git diff: %w", err)
	}

	return string(output), nil
}

func NewGit() VCS {
	return &git{}
}
