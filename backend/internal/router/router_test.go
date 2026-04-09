package router_test

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/router"
)

type mockAnalyzeUseCase struct{}

func (m *mockAnalyzeUseCase) Analyze(_ context.Context, _ multipart.File) (entity.AnalyzeResult, error) {
	return entity.AnalyzeResult{}, nil
}

func TestSetupRouterRegistersAnalyzeRoute(t *testing.T) {
	r := router.SetupRouter(&mockAnalyzeUseCase{}, 1024)

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing image, got %d", w.Code)
	}
}
