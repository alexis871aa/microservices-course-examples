package v1

import (
	"context"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/ufo/v1"
)

// AnalyzeSighting анализирует наблюдение НЛО
func (a *api) AnalyzeSighting(ctx context.Context, req *ufoV1.AnalyzeSightingRequest) (*ufoV1.AnalyzeSightingResponse, error) {
	result, err := a.ufoService.AnalyzeSighting(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &ufoV1.AnalyzeSightingResponse{
		AnalysisResult:  result.AnalysisResult,
		Classification:  result.Classification,
		ConfidenceScore: result.ConfidenceScore,
	}, nil
}
