package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// HTTP
	HTTP_ADDR string

	// Database
	DatabaseURL string

	// Kafka
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	// Cache
	CacheCapacity int
}

func Load() (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	// HTTP
	cfg.HTTP_ADDR = os.Getenv("HTTP_ADDR")
	if cfg.HTTP_ADDR == "" {
		return nil, fmt.Errorf("HTTP_ADDR is required")
	}

	// Database
	cfg.DatabaseURL = os.Getenv("DB_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}

	// Kafka
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}
	cfg.KafkaBrokers = strings.Split(kafkaBrokers, ",")

	cfg.KafkaTopic = os.Getenv("KAFKA_TOPIC")
	if cfg.KafkaTopic == "" {
		return nil, fmt.Errorf("KAFKA_TOPIC is required")
	}

	cfg.KafkaGroupID = os.Getenv("KAFKA_GROUP_ID")

	// Cache
	cacheCapacity := 1000
	if envVal := os.Getenv("CACHE_CAPACITY"); envVal != "" {
		if val, err := strconv.Atoi(envVal); err != nil {
			log.Printf("Invalid CACHE_CAPACITY '%s', using default: %d", envVal, cacheCapacity)
		} else if val <= 0 {
			log.Printf("CACHE_CAPACITY must be positive, using default: %d", cacheCapacity)
		} else {
			cacheCapacity = val
		}
	}
	cfg.CacheCapacity = cacheCapacity
	return cfg, nil
}
