package test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

type mockIAService struct {
	result entity.AnalyzeResult
	err    error
}

func (m *mockIAService) AnalyzeImage(_ context.Context, _ multipart.File) (entity.AnalyzeResult, error) {
	return m.result, m.err
}

type fakeMultipartFile struct {
	*bytes.Reader
}

func (f *fakeMultipartFile) Close() error {
	return nil
}

func newFakeMultipartFile(content string) multipart.File {
	return &fakeMultipartFile{Reader: bytes.NewReader([]byte(content))}
}

func TestUseCaseAnalyze_Success(t *testing.T) {
	uc := usecase.NewAnalyzeImageUseCase(&mockIAService{
		result: entity.AnalyzeResult{Tags: []entity.Tag{{Label: "Animal", Confidence: 0.9}}},
	})

	result, err := uc.Analyze(context.Background(), newFakeMultipartFile("image"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result.Tags) != 1 || result.Tags[0].Label != "Animal" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestUseCaseAnalyze_NilFile(t *testing.T) {
	uc := usecase.NewAnalyzeImageUseCase(&mockIAService{})

	_, err := uc.Analyze(context.Background(), nil)
	if apperrors.CodeOf(err) != apperrors.CodeInvalidRequest {
		t.Fatalf("expected INVALID_REQUEST, got %v", apperrors.CodeOf(err))
	}
}

func TestUseCaseAnalyze_UnknownErrorWrapped(t *testing.T) {
	uc := usecase.NewAnalyzeImageUseCase(&mockIAService{err: errors.New("unexpected")})

	_, err := uc.Analyze(context.Background(), newFakeMultipartFile("image"))
	if apperrors.CodeOf(err) != apperrors.CodeInternal {
		t.Fatalf("expected INTERNAL_ERROR, got %v", apperrors.CodeOf(err))
	}
}
