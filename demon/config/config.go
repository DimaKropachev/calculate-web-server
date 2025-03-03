package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	NumWorkers int
}

func LoadConfig() (*Config, error) {
	numWorker := os.Getenv("COMPUTING_POWER")
	if numWorker == "" {
		numWorker = "3"
	}

	n, err := strconv.Atoi(numWorker)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке конфигурации демона: %w", err)
	}

	return &Config{
		NumWorkers: n,
	}, nil
}
