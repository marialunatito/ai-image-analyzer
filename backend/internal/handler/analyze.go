package handler

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/usecase"
)

var allowedMIMETypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/jpg":  {},
	"image/webp": {},
}

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
