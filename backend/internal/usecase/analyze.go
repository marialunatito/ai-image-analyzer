package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
)

type AnalyzeImageUseCase interface {
	Analyze(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error)
}

type AnalyzeImageUseCaseImpl struct {
	service service.IAService
}

func NewAnalyzeImageUseCase(service service.IAService) AnalyzeImageUseCase {
	return &AnalyzeImageUseCaseImpl{service: service}
}

func (u *AnalyzeImageUseCaseImpl) Analyze(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error) {
	if file == nil {
		return entity.AnalyzeResult{}, apperrors.New(apperrors.CodeInvalidRequest, "el archivo de imagen es requerido")
	}

	result, err := u.service.AnalyzeImage(ctx, file)
	if err != nil {
		code := apperrors.CodeOf(err)
		switch code {
		case apperrors.CodeProviderError, apperrors.CodeProviderTimeout, apperrors.CodeInvalidImage:
			return entity.AnalyzeResult{}, err
		default:
			return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeInternal, fmt.Sprintf("error ejecutando analisis: %s", err.Error()), err)
		}
	}

	return result, nil
}
