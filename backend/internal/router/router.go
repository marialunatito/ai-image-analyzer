package router

import (
	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/handler"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler)
	return r
}
