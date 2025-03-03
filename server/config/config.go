package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type OperationTime struct {
	Addition       time.Duration
	Subtraction    time.Duration
	Multiplication time.Duration
	Division       time.Duration
}

type Config struct {
	Server   *ServerConfig
	Timeouts *OperationTime
}

func NewCongif() *Config {
	return &Config{
		Server:   &ServerConfig{},
		Timeouts: &OperationTime{},
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./server/config/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ошибка чтения конфигурации: %w", err))
	}

	config := NewCongif()

	err = viper.Unmarshal(config.Server)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	addStr := os.Getenv("TIME_ADDITION_MS")
	if addStr == "" {
		addStr = "100"
	}

	subStr := os.Getenv("TIME_SUBTRACTION_MS")
	if subStr == "" {
		subStr = "100"
	}

	multStr := os.Getenv("TIME_MULTIPLICATIONS_MS")
	if multStr == "" {
		multStr = "100"
	}

	divStr := os.Getenv("TIME_DIVISIONS_MS")
	if divStr == "" {
		divStr = "100"
	}

	add, err := strconv.Atoi(addStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	sub, err := strconv.Atoi(subStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	mult, err := strconv.Atoi(multStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	div, err := strconv.Atoi(divStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	times := &OperationTime{
		Addition:       time.Duration(add) * time.Millisecond,
		Subtraction:    time.Duration(sub) * time.Millisecond,
		Multiplication: time.Duration(mult) * time.Millisecond,
		Division:       time.Duration(div) * time.Millisecond,
	}

	config.Timeouts = times

	return config, nil
}
