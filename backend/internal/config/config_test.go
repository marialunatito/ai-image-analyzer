package config_test

import (
	"testing"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("GCV_API_KEY", "test-key")
	t.Setenv("PORT", "")
	t.Setenv("MAX_IMAGE_SIZE", "")
	t.Setenv("GCV_API_URL", "")
	t.Setenv("GCV_TIMEOUT_SECONDS", "")
	t.Setenv("GCV_MAX_RESULTS", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.MaxImageSize != 5*1024*1024 {
		t.Fatalf("unexpected max image size: %d", cfg.MaxImageSize)
	}
	if cfg.Vision.Timeout != 20*time.Second {
		t.Fatalf("expected default timeout 20s, got %v", cfg.Vision.Timeout)
	}
	if cfg.Vision.MaxResults != 10 {
		t.Fatalf("expected default max results 10, got %d", cfg.Vision.MaxResults)
	}
}

func TestLoadMissingAPIKey(t *testing.T) {
	t.Setenv("GCV_API_KEY", "")

	_, err := config.Load()
	if err == nil {
		t.Fatalf("expected error for missing API key")
	}
}

func TestLoadInvalidMaxImageSize(t *testing.T) {
	t.Setenv("GCV_API_KEY", "test-key")
	t.Setenv("MAX_IMAGE_SIZE", "0")

	_, err := config.Load()
	if err == nil {
		t.Fatalf("expected error for invalid max image size")
	}
}

func TestLoadSanitizesTimeoutAndMaxResults(t *testing.T) {
	t.Setenv("GCV_API_KEY", "test-key")
	t.Setenv("GCV_TIMEOUT_SECONDS", "-10")
	t.Setenv("GCV_MAX_RESULTS", "0")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.Vision.Timeout != 20*time.Second {
		t.Fatalf("expected timeout fallback 20s, got %v", cfg.Vision.Timeout)
	}
	if cfg.Vision.MaxResults != 10 {
		t.Fatalf("expected max results fallback 10, got %d", cfg.Vision.MaxResults)
	}
}
