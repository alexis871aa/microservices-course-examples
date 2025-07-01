package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
)

// AnalyzeSighting анализирует наблюдение НЛО через Analysis сервис
func (s *service) AnalyzeSighting(ctx context.Context, uuid string) (model.AnalysisResult, error) {
	// Создаем спан для вызова Analysis сервиса
	ctx, span := tracing.StartSpan(ctx, "ufo.call_analysis")
	defer span.End()

	// Вызываем Analysis сервис (клиент уже возвращает модель)
	return s.analysisClient.AnalyzeSighting(ctx, uuid)
}
