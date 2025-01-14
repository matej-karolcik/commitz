package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matej-karolcik/commitz/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dbPath := filepath.Join(filepath.Dir(file), "loaded.db")
	config, err := config.Load(dbPath)

	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dbPath)
	}()

	assert.Equal(t, config.Model, "llama3.2")
	assert.Equal(t, config.Temperature, 0.5)
}

func TestSaveAndLoad(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)

	cfg := &config.Config{
		Model:       "mistral",
		Temperature: 69.420,
	}

	dbPath := filepath.Join(filepath.Dir(file), "saved.db")
	err := cfg.Save(dbPath)
	assert.Nil(t, err)

	cfg, err = config.Load(dbPath)
	assert.Nil(t, err)

	defer func() {
		_ = os.RemoveAll(dbPath)
	}()

	assert.Equal(t, cfg.Model, "mistral")
	assert.Equal(t, cfg.Temperature, 69.420)
}
