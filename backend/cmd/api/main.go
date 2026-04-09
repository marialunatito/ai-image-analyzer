package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/api"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/config"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/router"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

type runnableServer interface {
	Run(addr ...string) error
}

var (
	loadConfig = config.Load
	newVisionService = func(apiKey, apiURL string, timeout time.Duration, maxResults int) service.IAService {
		return service.NewGoogleVisionService(apiKey, apiURL, timeout, maxResults)
	}
	newAnalyzeImageUseCase = usecase.NewAnalyzeImageUseCase
	setupHTTPRouter = func(analyzeUseCase usecase.AnalyzeImageUseCase, maxImageSize int64) runnableServer {
		return router.SetupRouter(analyzeUseCase, maxImageSize)
	}
)

// @title AI Image Analyzer API
// @version 1.0
// @description API para analizar imagenes usando un proveedor de IA.
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	visionService := newVisionService(
		cfg.Vision.APIKey,
		cfg.Vision.APIURL,
		cfg.Vision.Timeout,
		cfg.Vision.MaxResults,
	)
	analyzeUseCase := newAnalyzeImageUseCase(visionService)

	r := setupHTTPRouter(analyzeUseCase, cfg.MaxImageSize)
	if err := r.Run(":" + cfg.Port); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	return nil
}
