package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
)

type failingMultipartFile struct{}

type bytesMultipartFile struct {
	*bytes.Reader
}

func (f *failingMultipartFile) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func (f *failingMultipartFile) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, io.EOF
}

func (f *failingMultipartFile) Seek(_ int64, _ int) (int64, error) {
	return 0, nil
}

func (f *failingMultipartFile) Close() error {
	return nil
}

func (f *bytesMultipartFile) Close() error {
	return nil
}

type timeoutErr struct{}

func (e timeoutErr) Error() string { return "timeout" }
func (e timeoutErr) Timeout() bool { return true }

func TestAnalyzeImage_NilFile(t *testing.T) {
	svc := NewGoogleVisionService("key", "http://example.com", time.Second, 5)

	_, err := svc.AnalyzeImage(context.Background(), nil)
	if apperrors.CodeOf(err) != apperrors.CodeInvalidRequest {
		t.Fatalf("expected INVALID_REQUEST, got %v", apperrors.CodeOf(err))
	}
}

func TestAnalyzeImage_ReadError(t *testing.T) {
	svc := NewGoogleVisionService("key", "http://example.com", time.Second, 5)

	_, err := svc.AnalyzeImage(context.Background(), &failingMultipartFile{})
	if apperrors.CodeOf(err) != apperrors.CodeInternal {
		t.Fatalf("expected INTERNAL_ERROR, got %v", apperrors.CodeOf(err))
	}
}

func TestAnalyzeImage_EmptyContent(t *testing.T) {
	svc := NewGoogleVisionService("key", "http://example.com", time.Second, 5)
	file := &bytesMultipartFile{Reader: bytes.NewReader(nil)}

	_, err := svc.AnalyzeImage(context.Background(), file)
	if apperrors.CodeOf(err) != apperrors.CodeInvalidImage {
		t.Fatalf("expected INVALID_IMAGE, got %v", apperrors.CodeOf(err))
	}
}

func TestAnalyzeImage_InvalidEndpointURL(t *testing.T) {
	svc := NewGoogleVisionService("key", ":://bad url", time.Second, 5)
	file := &bytesMultipartFile{Reader: bytes.NewReader([]byte("img"))}

	_, err := svc.AnalyzeImage(context.Background(), file)
	if apperrors.CodeOf(err) != apperrors.CodeInternal {
		t.Fatalf("expected INTERNAL_ERROR, got %v", apperrors.CodeOf(err))
	}
}

func TestEndpoint(t *testing.T) {
	svc := NewGoogleVisionService("new-key", "https://vision.local/path?key=existing", time.Second, 5)
	googleSvc := svc.(*GoogleVisionService)

	if got := googleSvc.endpoint(); got != "https://vision.local/path?key=existing" {
		t.Fatalf("expected existing key to be preserved, got %q", got)
	}

	badSvc := &GoogleVisionService{apiURL: "%%%"}
	if got := badSvc.endpoint(); got != "%%%" {
		t.Fatalf("expected raw apiURL on parse error, got %q", got)
	}
}

func TestParseVisionResponse(t *testing.T) {
	_, err := parseVisionResponse([]byte("{"))
	if apperrors.CodeOf(err) != apperrors.CodeProviderError {
		t.Fatalf("expected PROVIDER_ERROR for invalid json, got %v", apperrors.CodeOf(err))
	}

	result, err := parseVisionResponse([]byte(`{"responses":[]}`))
	if err != nil {
		t.Fatalf("expected nil error for empty responses, got %v", err)
	}
	if len(result.Tags) != 0 {
		t.Fatalf("expected no tags, got %+v", result.Tags)
	}

	_, err = parseVisionResponse([]byte(`{"responses":[{"error":{"message":"provider failed"}}]}`))
	if apperrors.CodeOf(err) != apperrors.CodeProviderError {
		t.Fatalf("expected PROVIDER_ERROR for provider message, got %v", apperrors.CodeOf(err))
	}

	result, err = parseVisionResponse([]byte(`{"responses":[{"labelAnnotations":[{"description":"","score":0.1},{"description":"Dog","score":0.99}]}]}`))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result.Tags) != 1 || result.Tags[0].Label != "Dog" {
		t.Fatalf("unexpected tags: %+v", result.Tags)
	}
}

func TestErrorsIsTimeout(t *testing.T) {
	if errorsIsTimeout(nil, nil) {
		t.Fatal("expected nil error to not be timeout")
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	if !errorsIsTimeout(errors.New("any"), ctx) {
		t.Fatal("expected context deadline exceeded to be timeout")
	}

	if !errorsIsTimeout(timeoutErr{}, context.Background()) {
		t.Fatal("expected Timeout() error to be timeout")
	}
}

var _ multipart.File = (*failingMultipartFile)(nil)