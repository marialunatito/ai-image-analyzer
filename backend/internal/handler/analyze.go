package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnalyzeHandler(c *gin.Context) {
	// TODO: call the service or use case responsible for analyzing the image and return the result
	c.JSON(http.StatusOK, gin.H{
		"tags": []map[string]interface{}{
			{"label": "Gato", "confidence": 0.98},
			{"label": "Mascota", "confidence": 0.87},
			{"label": "Animal", "confidence": 0.92},
		},
	})

}