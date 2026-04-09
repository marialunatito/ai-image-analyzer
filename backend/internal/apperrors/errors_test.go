package apperrors_test

import (
	"errors"
	"testing"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
)

func TestNewAndCodeOf(t *testing.T) {
	err := apperrors.New(apperrors.CodeInvalidRequest, "bad request")
	if apperrors.CodeOf(err) != apperrors.CodeInvalidRequest {
		t.Fatalf("expected %s, got %s", apperrors.CodeInvalidRequest, apperrors.CodeOf(err))
	}
}

func TestWrapAndUnwrap(t *testing.T) {
	root := errors.New("root")
	err := apperrors.Wrap(apperrors.CodeProviderError, "provider failed", root)

	if !errors.Is(err, root) {
		t.Fatalf("expected wrapped error")
	}
	if apperrors.CodeOf(err) != apperrors.CodeProviderError {
		t.Fatalf("expected %s, got %s", apperrors.CodeProviderError, apperrors.CodeOf(err))
	}
}

func TestCodeOfUnknownErrorReturnsInternal(t *testing.T) {
	err := errors.New("unknown")
	if apperrors.CodeOf(err) != apperrors.CodeInternal {
		t.Fatalf("expected %s, got %s", apperrors.CodeInternal, apperrors.CodeOf(err))
	}
}

func TestNilAppErrorMethods(t *testing.T) {
	var e *apperrors.AppError
	if got := e.Error(); got != "" {
		t.Fatalf("expected empty error string, got %q", got)
	}
	if e.Unwrap() != nil {
		t.Fatalf("expected nil unwrap")
	}
}
