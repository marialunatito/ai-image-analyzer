package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/apperrors"
	"github.com/marialuna/prueba_tecnica/ai-image-analyzer/backend/internal/entity"
)

type IAService interface {
	AnalyzeImage(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error)
}

type GoogleVisionService struct {
	client     *http.Client
	apiKey     string
	apiURL     string
	maxResults int
}

func NewGoogleVisionService(apiKey, apiURL string, timeout time.Duration, maxResults int) IAService {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	if maxResults <= 0 {
		maxResults = 10
	}

	return &GoogleVisionService{
		client:     &http.Client{Timeout: timeout},
		apiKey:     strings.TrimSpace(apiKey),
		apiURL:     strings.TrimSpace(apiURL),
		maxResults: maxResults,
	}
}

func (s *GoogleVisionService) AnalyzeImage(ctx context.Context, file multipart.File) (entity.AnalyzeResult, error) {
	if file == nil {
		return entity.AnalyzeResult{}, apperrors.New(apperrors.CodeInvalidRequest, "el archivo de imagen es requerido")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeInternal, "no se pudo leer la imagen", err)
	}

	if len(content) == 0 {
		return entity.AnalyzeResult{}, apperrors.New(apperrors.CodeInvalidImage, "la imagen esta vacia")
	}

	body, err := s.buildRequestBody(content)
	if err != nil {
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeInternal, "no se pudo serializar la solicitud de vision", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint(), bytes.NewReader(body))
	if err != nil {
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeInternal, "no se pudo crear la solicitud HTTP", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		if errorsIsTimeout(err, ctx) {
			return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeProviderTimeout, "el proveedor de IA excedio el tiempo de espera", err)
		}
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeProviderError, "fallo la comunicacion con el proveedor de IA", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeProviderError, "no se pudo leer la respuesta del proveedor", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return entity.AnalyzeResult{}, apperrors.New(apperrors.CodeProviderError, fmt.Sprintf("proveedor IA respondio con estado %d", resp.StatusCode))
	}

	return parseVisionResponse(responseBody)
}

type visionRequest struct {
	Requests []visionAnnotateRequest `json:"requests"`
}

type visionAnnotateRequest struct {
	Image    visionImage     `json:"image"`
	Features []visionFeature `json:"features"`
}

type visionImage struct {
	Content string `json:"content"`
}

type visionFeature struct {
	Type       string `json:"type"`
	MaxResults int    `json:"maxResults"`
}

type visionResponse struct {
	Responses []visionAnnotateResponse `json:"responses"`
}

type visionAnnotateResponse struct {
	LabelAnnotations []visionLabelAnnotation `json:"labelAnnotations"`
	Error            *visionProviderError    `json:"error,omitempty"`
}

type visionLabelAnnotation struct {
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type visionProviderError struct {
	Message string `json:"message"`
}

func (s *GoogleVisionService) buildRequestBody(content []byte) ([]byte, error) {
	payload := visionRequest{
		Requests: []visionAnnotateRequest{
			{
				Image: visionImage{Content: base64.StdEncoding.EncodeToString(content)},
				Features: []visionFeature{
					{Type: "LABEL_DETECTION", MaxResults: s.maxResults},
				},
			},
		},
	}

	return json.Marshal(payload)
}

func (s *GoogleVisionService) endpoint() string {
	parsed, err := url.Parse(s.apiURL)
	if err != nil {
		return s.apiURL
	}
	query := parsed.Query()
	if query.Get("key") == "" && s.apiKey != "" {
		query.Set("key", s.apiKey)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func parseVisionResponse(raw []byte) (entity.AnalyzeResult, error) {
	var response visionResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return entity.AnalyzeResult{}, apperrors.Wrap(apperrors.CodeProviderError, "respuesta del proveedor no valida", err)
	}

	if len(response.Responses) == 0 {
		return entity.AnalyzeResult{Tags: []entity.Tag{}}, nil
	}

	first := response.Responses[0]
	if first.Error != nil && strings.TrimSpace(first.Error.Message) != "" {
		return entity.AnalyzeResult{}, apperrors.New(apperrors.CodeProviderError, first.Error.Message)
	}

	tags := make([]entity.Tag, 0, len(first.LabelAnnotations))
	for _, label := range first.LabelAnnotations {
		if strings.TrimSpace(label.Description) == "" {
			continue
		}
		tags = append(tags, entity.Tag{
			Label:      label.Description,
			Confidence: label.Score,
		})
	}

	return entity.AnalyzeResult{Tags: tags}, nil
}

func errorsIsTimeout(err error, ctx context.Context) bool {
	if err == nil {
		return false
	}

	if ctx != nil && ctx.Err() == context.DeadlineExceeded {
		return true
	}

	timeoutErr, ok := err.(interface{ Timeout() bool })
	return ok && timeoutErr.Timeout()
}
