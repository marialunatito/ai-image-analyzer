package usecase

import (
	"context"
	"mime/multipart"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
)

type AnalyzeImageUseCase interface {
	Analyze(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error)
}
