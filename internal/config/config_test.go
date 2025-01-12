package config_test

import (
	"testing"

	"github.com/matej-karolcik/commitz/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	config, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	assert.Equal(t, config.Model, "llama3.2")
	assert.Equal(t, config.Temperature, 0.5)
}

func TestLoadEmpty(t *testing.T) {
	config, err := config.Load(".empty.yaml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	assert.Equal(t, "llama3.2", config.Model)
	assert.Equal(t, 0.5, config.Temperature)
}
