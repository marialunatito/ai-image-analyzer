package main

import (
	"log"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/config"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/router"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	visionService := service.NewGoogleVisionService(
		cfg.Vision.APIKey,
		cfg.Vision.APIURL,
		cfg.Vision.Timeout,
		cfg.Vision.MaxResults,
	)
	analyzeUseCase := usecase.NewAnalyzeImageUseCase(visionService)

	r := router.SetupRouter(analyzeUseCase, cfg.MaxImageSize)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
