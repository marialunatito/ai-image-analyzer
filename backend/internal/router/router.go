package router

import (
	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/handler"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

func SetupRouter(analyzeUseCase usecase.AnalyzeImageUseCase, maxImageSize int64) *gin.Engine {
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler(analyzeUseCase, maxImageSize))
	return r
}
