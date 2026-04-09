package main

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/config"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

type fakeServer struct {
	err          error
	receivedAddr string
}

func (s *fakeServer) Run(addr ...string) error {
	if len(addr) > 0 {
		s.receivedAddr = addr[0]
	}
	return s.err
}

type fakeIAService struct{}

func (f *fakeIAService) AnalyzeImage(_ context.Context, _ multipart.File) (entity.AnalyzeResult, error) {
	return entity.AnalyzeResult{}, nil
}

func TestRun_LoadConfigError(t *testing.T) {
	origLoadConfig := loadConfig
	defer func() { loadConfig = origLoadConfig }()

	loadConfig = func() (config.Config, error) {
		return config.Config{}, errors.New("boom")
	}

	err := run()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRun_ServerRunError(t *testing.T) {
	origLoadConfig := loadConfig
	origNewVisionService := newVisionService
	origNewAnalyzeImageUseCase := newAnalyzeImageUseCase
	origSetupHTTPRouter := setupHTTPRouter
	defer func() {
		loadConfig = origLoadConfig
		newVisionService = origNewVisionService
		newAnalyzeImageUseCase = origNewAnalyzeImageUseCase
		setupHTTPRouter = origSetupHTTPRouter
	}()

	loadConfig = func() (config.Config, error) {
		return config.Config{
			Port:         "9999",
			MaxImageSize: 1024,
			Vision: config.VisionConfig{
				APIKey:     "k",
				APIURL:     "http://example.com",
				Timeout:    time.Second,
				MaxResults: 1,
			},
		}, nil
	}

	newVisionService = func(apiKey, apiURL string, timeout time.Duration, maxResults int) service.IAService {
		return &fakeIAService{}
	}

	newAnalyzeImageUseCase = func(svc service.IAService) usecase.AnalyzeImageUseCase {
		return usecase.NewAnalyzeImageUseCase(svc)
	}

	server := &fakeServer{err: errors.New("run failed")}
	setupHTTPRouter = func(_ usecase.AnalyzeImageUseCase, _ int64) runnableServer {
		return server
	}

	err := run()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if server.receivedAddr != ":9999" {
		t.Fatalf("expected addr :9999, got %q", server.receivedAddr)
	}
}

func TestRun_Success(t *testing.T) {
	origLoadConfig := loadConfig
	origNewVisionService := newVisionService
	origNewAnalyzeImageUseCase := newAnalyzeImageUseCase
	origSetupHTTPRouter := setupHTTPRouter
	defer func() {
		loadConfig = origLoadConfig
		newVisionService = origNewVisionService
		newAnalyzeImageUseCase = origNewAnalyzeImageUseCase
		setupHTTPRouter = origSetupHTTPRouter
	}()

	loadConfig = func() (config.Config, error) {
		return config.Config{
			Port:         "8081",
			MaxImageSize: 2048,
			Vision: config.VisionConfig{
				APIKey:     "k",
				APIURL:     "http://example.com",
				Timeout:    time.Second,
				MaxResults: 2,
			},
		}, nil
	}

	newVisionService = func(apiKey, apiURL string, timeout time.Duration, maxResults int) service.IAService {
		return &fakeIAService{}
	}

	newAnalyzeImageUseCase = func(svc service.IAService) usecase.AnalyzeImageUseCase {
		return usecase.NewAnalyzeImageUseCase(svc)
	}

	server := &fakeServer{}
	setupHTTPRouter = func(_ usecase.AnalyzeImageUseCase, _ int64) runnableServer {
		return server
	}

	if err := run(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if server.receivedAddr != ":8081" {
		t.Fatalf("expected addr :8081, got %q", server.receivedAddr)
	}
}
