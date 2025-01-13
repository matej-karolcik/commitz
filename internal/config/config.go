package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	Model       string  `default:"llama3.2"`
	Temperature float64 `default:"0.5"`
}

func Load(paths ...string) (*Config, error) {
	db, err := openDB(paths...)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close badger", "error", err)
		}
	}()

	var (
		model       string
		temperature float64
	)

	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("model"))
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf("failed to get model: %w", err)
		} else if err == nil {
			if err := item.Value(func(val []byte) error {
				model = string(val)
				return nil
			}); err != nil {
				return fmt.Errorf("failed to get model: %w", err)
			}
		}

		item, err = txn.Get([]byte("temperature"))
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf("failed to get temperature: %w", err)
		} else if err == nil {
			if err := item.Value(func(val []byte) error {
				temperature, err = strconv.ParseFloat(string(val), 64)
				if err != nil {
					return fmt.Errorf("failed to parse temperature: %w", err)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("failed to get temperature: %w", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	config := Config{
		Model:       model,
		Temperature: temperature,
	}

	defaults.SetDefaults(&config)

	return &config, nil
}

func (c *Config) Save(paths ...string) error {
	db, err := openDB(paths...)
	if err != nil {
		return fmt.Errorf("failed to open badger: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close badger", "error", err)
		}
	}()

	if err := db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte("model"), []byte(c.Model)); err != nil {
			return fmt.Errorf("failed to set model: %w", err)
		}

		if err := txn.Set([]byte("temperature"), []byte(strconv.FormatFloat(c.Temperature, 'f', -1, 64))); err != nil {
			return fmt.Errorf("failed to set temperature: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func openDB(paths ...string) (*badger.DB, error) {
	var path string
	if len(paths) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		path = filepath.Join(homeDir, ".commitz")
	} else {
		path = paths[0]
	}

	db, err := badger.Open(badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}

	return db, nil
}
