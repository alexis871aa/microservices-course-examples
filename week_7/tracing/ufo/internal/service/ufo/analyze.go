package ufo

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
)

// AnalyzeSighting анализирует наблюдение НЛО через Analysis сервис
func (s *service) AnalyzeSighting(ctx context.Context, uuid string) (model.AnalysisResult, error) {
	// Создаем спан для вызова Analysis сервиса
	ctx, span := tracing.StartSpan(ctx, "ufo.call_analysis",
		trace.WithAttributes(
			attribute.String("ufo.uuid", uuid),
		),
	)
	defer span.End()

	// Вызываем Analysis сервис (клиент уже возвращает модель)
	result, err := s.analysisClient.AnalyzeSighting(ctx, uuid)
	if err != nil {
		span.RecordError(err)
		return model.AnalysisResult{}, err
	}

	// Добавляем атрибуты результата
	span.SetAttributes(
		attribute.String("analysis.classification", result.Classification),
		attribute.Float64("analysis.confidence", float64(result.ConfidenceScore)),
		attribute.String("analysis.result", result.AnalysisResult),
	)

	return result, nil
}
