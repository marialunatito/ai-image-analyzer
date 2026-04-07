package service

import (
	"context"
	"mime/multipart"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
)

type IAService interface {
	AnalyzeImage(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error)
}
