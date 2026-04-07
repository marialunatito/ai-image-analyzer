package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/handler"
)

func TestAnalyzeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/analyze", handler.AnalyzeHandler)

	req := httptest.NewRequest("POST", "/api/analyze", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"tags":[{"confidence":0.98,"label":"Gato"},{"confidence":0.87,"label":"Mascota"},{"confidence":0.92,"label":"Animal"}]}`

	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}
