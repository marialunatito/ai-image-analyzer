package handler

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

type AnalyzeResponse = entity.AnalyzeResult

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var allowedMIMETypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/jpg":  {},
	"image/webp": {},
}

// AnalyzeHandler godoc
// @Summary Analizar imagen
// @Description Recibe una imagen y devuelve etiquetas de contenido con su confianza.
// @Tags analyze
// @Accept mpfd
// @Produce json
// @Param image formData file true "Archivo de imagen"
// @Success 200 {object} AnalyzeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 415 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Failure 504 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/analyze [post]
func AnalyzeHandler(analyzeUseCase usecase.AnalyzeImageUseCase, maxImageSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			writeError(c, apperrors.New(apperrors.CodeInvalidRequest, "el campo image es requerido"))
			return
		}

		if maxImageSize > 0 && fileHeader.Size > maxImageSize {
			writeError(c, apperrors.New(apperrors.CodePayloadLarge, "la imagen supera el tamano maximo permitido"))
			return
		}

		if !isAllowedImageType(fileHeader.Header.Get("Content-Type"), fileHeader.Filename) {
			writeError(c, apperrors.New(apperrors.CodeInvalidImage, "formato de imagen no permitido"))
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			writeError(c, apperrors.Wrap(apperrors.CodeInvalidRequest, "no se pudo abrir la imagen enviada", err))
			return
		}
		defer file.Close()

		result, err := analyzeUseCase.Analyze(c.Request.Context(), file)
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func writeError(c *gin.Context, err error) {
	code := apperrors.CodeOf(err)
	status := httpStatusForCode(code)

	message := "error interno al procesar la imagen"
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) && strings.TrimSpace(appErr.Message) != "" {
		message = appErr.Message
	}

	c.JSON(status, gin.H{
		"message": message,
		"code":    string(code),
	})
}

func httpStatusForCode(code apperrors.Code) int {
	switch code {
	case apperrors.CodeInvalidRequest:
		return http.StatusBadRequest
	case apperrors.CodeInvalidImage:
		return http.StatusUnsupportedMediaType
	case apperrors.CodePayloadLarge:
		return http.StatusRequestEntityTooLarge
	case apperrors.CodeProviderTimeout:
		return http.StatusGatewayTimeout
	case apperrors.CodeProviderError:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func isAllowedImageType(contentType, fileName string) bool {
	normalized := normalizeContentType(contentType)
	if normalized != "" {
		_, ok := allowedMIMETypes[normalized]
		return ok
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return false
	}
	detected := normalizeContentType(mime.TypeByExtension(ext))
	_, ok := allowedMIMETypes[detected]
	return ok
}

func normalizeContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(contentType))
	if err != nil {
		return strings.ToLower(strings.TrimSpace(contentType))
	}
	return strings.ToLower(mediaType)
}
