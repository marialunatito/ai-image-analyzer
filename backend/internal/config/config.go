package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort         = "8080"
	defaultMaxImageSize = int64(5 * 1024 * 1024)
	defaultVisionURL    = "https://vision.googleapis.com/v1/images:annotate"
	defaultTimeoutSec   = 20
	defaultMaxResults   = 10
)

type Config struct {
	Port         string
	MaxImageSize int64
	Vision       VisionConfig
}

type VisionConfig struct {
	APIKey     string
	APIURL     string
	Timeout    time.Duration
	MaxResults int
}

func Load() (Config, error) {
	cfg := Config{
		Port:         readString("PORT", defaultPort),
		MaxImageSize: readInt64("MAX_IMAGE_SIZE", defaultMaxImageSize),
		Vision: VisionConfig{
			APIKey:     readString("GCV_API_KEY", ""),
			APIURL:     readString("GCV_API_URL", defaultVisionURL),
			Timeout:    time.Duration(readInt("GCV_TIMEOUT_SECONDS", defaultTimeoutSec)) * time.Second,
			MaxResults: readInt("GCV_MAX_RESULTS", defaultMaxResults),
		},
	}

	if strings.TrimSpace(cfg.Vision.APIKey) == "" {
		return Config{}, fmt.Errorf("missing GCV_API_KEY")
	}

	if cfg.MaxImageSize <= 0 {
		return Config{}, fmt.Errorf("MAX_IMAGE_SIZE must be greater than 0")
	}

	if cfg.Vision.MaxResults <= 0 {
		cfg.Vision.MaxResults = defaultMaxResults
	}

	if cfg.Vision.Timeout <= 0 {
		cfg.Vision.Timeout = time.Duration(defaultTimeoutSec) * time.Second
	}

	return cfg, nil
}

func readString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func readInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func readInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
