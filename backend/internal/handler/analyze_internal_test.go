package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
)

func TestHTTPStatusForCode(t *testing.T) {
	tests := []struct {
		name string
		code apperrors.Code
		want int
	}{
		{name: "invalid request", code: apperrors.CodeInvalidRequest, want: http.StatusBadRequest},
		{name: "invalid image", code: apperrors.CodeInvalidImage, want: http.StatusUnsupportedMediaType},
		{name: "payload too large", code: apperrors.CodePayloadLarge, want: http.StatusRequestEntityTooLarge},
		{name: "provider timeout", code: apperrors.CodeProviderTimeout, want: http.StatusGatewayTimeout},
		{name: "provider error", code: apperrors.CodeProviderError, want: http.StatusBadGateway},
		{name: "default internal", code: apperrors.CodeInternal, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpStatusForCode(tt.code); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestNormalizeContentType(t *testing.T) {
	if got := normalizeContentType(" image/jpeg; charset=binary "); got != "image/jpeg" {
		t.Fatalf("expected image/jpeg, got %q", got)
	}

	if got := normalizeContentType(" IMAGE/PNG "); got != "image/png" {
		t.Fatalf("expected image/png, got %q", got)
	}
}

func TestIsAllowedImageType(t *testing.T) {
	if !isAllowedImageType("image/png", "file.any") {
		t.Fatal("expected png content type to be allowed")
	}

	if !isAllowedImageType("", "photo.JPG") {
		t.Fatal("expected jpg extension fallback to be allowed")
	}

	if isAllowedImageType("", "file.txt") {
		t.Fatal("expected txt extension to be rejected")
	}

	if isAllowedImageType("", "no-extension") {
		t.Fatal("expected missing extension to be rejected")
	}
}

func TestWriteError_DefaultMessageForNonAppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	writeError(c, errors.New("unexpected"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if payload["code"] != string(apperrors.CodeInternal) {
		t.Fatalf("expected INTERNAL_ERROR, got %q", payload["code"])
	}

	if payload["message"] != "error interno al procesar la imagen" {
		t.Fatalf("unexpected message: %q", payload["message"])
	}
}

func TestWriteError_UsesDefaultWhenAppErrorMessageIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	writeError(c, apperrors.New(apperrors.CodeProviderError, "  "))

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if payload["code"] != string(apperrors.CodeProviderError) {
		t.Fatalf("expected PROVIDER_ERROR, got %q", payload["code"])
	}

	if payload["message"] != "error interno al procesar la imagen" {
		t.Fatalf("unexpected message: %q", payload["message"])
	}
}