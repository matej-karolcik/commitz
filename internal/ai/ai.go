package ai

import (
	"context"
)

type AI interface {
	CommitMessage(ctx context.Context, diff string) (string, error)
	ReadmeFile(ctx context.Context, filemap map[string]string) (string, error)
}
