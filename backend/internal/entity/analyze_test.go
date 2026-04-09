package entity_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
)

func TestAnalyzeResultJSONTags(t *testing.T) {
	result := entity.AnalyzeResult{
		Tags: []entity.Tag{{Label: "Dog", Confidence: 0.98}},
	}

	b, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	jsonStr := string(b)
	if !strings.Contains(jsonStr, "\"tags\"") {
		t.Fatalf("expected json to contain tags field, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, "\"label\"") {
		t.Fatalf("expected json to contain label field, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, "\"confidence\"") {
		t.Fatalf("expected json to contain confidence field, got %s", jsonStr)
	}
}
