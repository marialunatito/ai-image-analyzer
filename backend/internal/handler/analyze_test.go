package handler_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/handler"
)

type mockAnalyzeUseCase struct {
	result entity.AnalyzeResult
	err    error
}

func (m *mockAnalyzeUseCase) Analyze(_ context.Context, _ multipart.File) (entity.AnalyzeResult, error) {
	return m.result, m.err
}

func TestAnalyzeHandler_MissingImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler(&mockAnalyzeUseCase{}, 5*1024*1024))

	req := httptest.NewRequest("POST", "/api/analyze", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAnalyzeHandler_InvalidImageType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler(&mockAnalyzeUseCase{}, 5*1024*1024))

	body, contentType := buildMultipartBody(t, "image", "file.txt", "text/plain", []byte("not-an-image"))
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status 415, got %d", w.Code)
	}
}

func TestAnalyzeHandler_FileTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler(&mockAnalyzeUseCase{}, 2))

	body, contentType := buildMultipartBody(t, "image", "image.png", "image/png", []byte("1234"))
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", w.Code)
	}
}

func TestAnalyzeHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockUC := &mockAnalyzeUseCase{
		result: entity.AnalyzeResult{Tags: []entity.Tag{{Label: "Perro", Confidence: 0.99}}},
	}
	r.POST("/api/analyze", handler.AnalyzeHandler(mockUC, 5*1024*1024))

	body, contentType := buildMultipartBody(t, "image", "dog.png", "image/png", []byte("png-content"))
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	expected := `{"tags":[{"label":"Perro","confidence":0.99}]}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func buildMultipartBody(t *testing.T, fieldName, fileName, contentType string, payload []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+fileName+`"`)
	partHeader.Set("Content-Type", contentType)

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("unable to create multipart part: %v", err)
	}

	if _, err := part.Write(payload); err != nil {
		t.Fatalf("unable to write multipart payload: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("unable to close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}