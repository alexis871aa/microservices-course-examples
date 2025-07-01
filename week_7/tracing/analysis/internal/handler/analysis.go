package handler

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/analysis/internal/service"
	analysisV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/analysis/v1"
)

// AnalysisHandler - gRPC хендлер для анализа
type AnalysisHandler struct {
	analysisV1.UnimplementedAnalysisServiceServer
	service *service.AnalysisService
}

// NewAnalysisHandler создает новый хендлер
func NewAnalysisHandler(service *service.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		service: service,
	}
}

// AnalyzeSighting обрабатывает запрос на анализ наблюдения
func (h *AnalysisHandler) AnalyzeSighting(ctx context.Context, req *analysisV1.AnalyzeSightingRequest) (*analysisV1.AnalyzeSightingResponse, error) {
	analysisResult, classification, confidence, err := h.service.AnalyzeSighting(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &analysisV1.AnalyzeSightingResponse{
		AnalysisResult:  analysisResult,
		Classification:  classification,
		ConfidenceScore: confidence,
	}, nil
}
