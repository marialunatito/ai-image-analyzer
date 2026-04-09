package service_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
)

type serviceMultipartFile struct {
	*bytes.Reader
}

func (f *serviceMultipartFile) Close() error {
	return nil
}

func newServiceMultipartFile(content []byte) multipart.File {
	return &serviceMultipartFile{Reader: bytes.NewReader(content)}
}

func TestGoogleVisionService_AnalyzeImageSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Fatalf("expected API key in query, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"responses":[{"labelAnnotations":[{"description":"Dog","score":0.98}]}]}`))
	}))
	defer server.Close()

	svc := service.NewGoogleVisionService("test-key", server.URL, 2*time.Second, 5)
	result, err := svc.AnalyzeImage(context.Background(), newServiceMultipartFile([]byte("fake-image")))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result.Tags) != 1 || result.Tags[0].Label != "Dog" {
		t.Fatalf("unexpected tags: %+v", result.Tags)
	}
}

func TestGoogleVisionService_ProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"quota exceeded"}}`))
	}))
	defer server.Close()

	svc := service.NewGoogleVisionService("test-key", server.URL, 2*time.Second, 5)
	_, err := svc.AnalyzeImage(context.Background(), newServiceMultipartFile([]byte("fake-image")))
	if apperrors.CodeOf(err) != apperrors.CodeProviderError {
		t.Fatalf("expected PROVIDER_ERROR, got %v", apperrors.CodeOf(err))
	}
}

func TestGoogleVisionService_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"responses":[]}`))
	}))
	defer server.Close()

	svc := service.NewGoogleVisionService("test-key", server.URL, 10*time.Millisecond, 5)
	_, err := svc.AnalyzeImage(context.Background(), newServiceMultipartFile([]byte("fake-image")))
	if apperrors.CodeOf(err) != apperrors.CodeProviderTimeout {
		t.Fatalf("expected PROVIDER_TIMEOUT, got %v", apperrors.CodeOf(err))
	}
}