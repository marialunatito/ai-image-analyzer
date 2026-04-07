package main

import (
	"os"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := router.SetupRouter()
	r.Run(":" + port)
}
