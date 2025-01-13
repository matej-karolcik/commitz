package ai

import (
	"context"
)

type AI interface {
	CommitMessage(ctx context.Context, diff string) (string, error)
}
