package router

import (
	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/handler"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(analyzeUseCase usecase.AnalyzeImageUseCase, maxImageSize int64) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/analyze", handler.AnalyzeHandler(analyzeUseCase, maxImageSize))
	return r
}
