package config_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matej-karolcik/commitz/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	config, err := config.Load(filepath.Join(filepath.Dir(file), "loaded.db"))

	assert.Nil(t, err)
	assert.Equal(t, config.Model, "llama3.2")
	assert.Equal(t, config.Temperature, 0.5)
}

func TestSaveAndLoad(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)

	cfg := &config.Config{
		Model:       "mistral",
		Temperature: 69.420,
	}

	err := cfg.Save(filepath.Join(filepath.Dir(file), "saved.db"))
	assert.Nil(t, err)

	cfg, err = config.Load(filepath.Join(filepath.Dir(file), "saved.db"))
	assert.Nil(t, err)
	assert.Equal(t, cfg.Model, "mistral")
	assert.Equal(t, cfg.Temperature, 69.420)
}
