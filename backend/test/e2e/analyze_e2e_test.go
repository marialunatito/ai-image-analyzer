package e2e_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/router"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/service"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

type analyzeResponse struct {
	Tags []struct {
		Label      string  `json:"label"`
		Confidence float64 `json:"confidence"`
	} `json:"tags"`
}

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func TestAnalyzeE2E_Success(t *testing.T) {
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Fatalf("expected api key test-key, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"responses":[{"labelAnnotations":[{"description":"Dog","score":0.98},{"description":"Park","score":0.88}]}]}`))
	}))
	defer provider.Close()

	svc := service.NewGoogleVisionService("test-key", provider.URL, 2*time.Second, 5)
	uc := usecase.NewAnalyzeImageUseCase(svc)
	r := router.SetupRouter(uc, 5*1024*1024)

	body, contentType := buildMultipartBody(t, "image", "dog.png", "image/png", []byte("fake-png-content"))
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", w.Code, w.Body.String())
	}

	var response analyzeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}

	if len(response.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(response.Tags))
	}
	if response.Tags[0].Label != "Dog" {
		t.Fatalf("expected first tag Dog, got %q", response.Tags[0].Label)
	}
}

func TestAnalyzeE2E_ProviderError(t *testing.T) {
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"quota exceeded"}}`))
	}))
	defer provider.Close()

	svc := service.NewGoogleVisionService("test-key", provider.URL, 2*time.Second, 5)
	uc := usecase.NewAnalyzeImageUseCase(svc)
	r := router.SetupRouter(uc, 5*1024*1024)

	body, contentType := buildMultipartBody(t, "image", "dog.png", "image/png", []byte("fake-png-content"))
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d with body %s", w.Code, w.Body.String())
	}

	var response errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON error response, got error: %v", err)
	}
	if response.Code != "PROVIDER_ERROR" {
		t.Fatalf("expected error code PROVIDER_ERROR, got %q", response.Code)
	}
}

func buildMultipartBody(t *testing.T, fieldName, fileName, contentType string, payload []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+fileName+`"`)
	headers.Set("Content-Type", contentType)

	part, err := writer.CreatePart(headers)
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
