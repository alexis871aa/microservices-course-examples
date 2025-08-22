package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
	classificationV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/classification/v1"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/ufo/v1"
)

// AnalysisService - сервис для анализа наблюдений НЛО
type AnalysisService struct {
	ufoClient            ufoV1.UFOServiceClient
	classificationClient classificationV1.ClassificationServiceClient
}

// NewAnalysisService создает новый сервис анализа
func NewAnalysisService() (*AnalysisService, error) {
	// Подключение к UFO сервису
	ufoConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("ufo-service")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to UFO service: %w", err)
	}

	// Подключение к Classification сервису
	classificationConn, err := grpc.NewClient(
		"localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("classification-service")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Classification service: %w", err)
	}

	return &AnalysisService{
		ufoClient:            ufoV1.NewUFOServiceClient(ufoConn),
		classificationClient: classificationV1.NewClassificationServiceClient(classificationConn),
	}, nil
}

// AnalyzeSighting анализирует наблюдение НЛО
func (s *AnalysisService) AnalyzeSighting(ctx context.Context, uuid string) (string, string, float32, error) {
	// Спан для получения данных о наблюдении
	ctx, getSightingSpan := tracing.StartSpan(ctx, "ufo.get_sighting",
		trace.WithAttributes(
			attribute.String("ufo.uuid", uuid),
		),
	)

	// Получаем данные о наблюдении из UFO сервиса
	sighting, err := s.ufoClient.Get(ctx, &ufoV1.GetRequest{Uuid: uuid})
	if err != nil {
		getSightingSpan.RecordError(err)
		getSightingSpan.End()
		return "", "", 0, fmt.Errorf("failed to get sighting: %w", err)
	}

	// Добавляем атрибуты после получения данных
	getSightingSpan.SetAttributes(
		attribute.String("sighting.description", sighting.Sighting.Info.Description),
		attribute.String("sighting.color", sighting.Sighting.Info.Color.GetValue()),
		attribute.Int("sighting.duration_seconds", int(sighting.Sighting.Info.DurationSeconds.GetValue())),
	)
	getSightingSpan.End()

	// Спан для классификации
	ctx, classifySpan := tracing.StartSpan(ctx, "analysis.classify",
		trace.WithAttributes(
			attribute.String("ufo.uuid", uuid),
			attribute.String("sighting.description", sighting.Sighting.Info.Description),
		),
	)
	defer classifySpan.End()

	// Отправляем данные в Classification сервис
	classification, err := s.classificationClient.ClassifyObject(ctx, &classificationV1.ClassifyObjectRequest{
		Description:     sighting.Sighting.Info.Description,
		Color:           sighting.Sighting.Info.Color.GetValue(),
		DurationSeconds: sighting.Sighting.Info.DurationSeconds.GetValue(),
	})
	if err != nil {
		classifySpan.RecordError(err)
		return "", "", 0, fmt.Errorf("failed to classify object: %w", err)
	}

	// Добавляем атрибуты результата классификации
	classifySpan.SetAttributes(
		attribute.String("analysis.classification", classification.ObjectType),
		attribute.Float64("analysis.confidence", float64(classification.Confidence)),
		attribute.String("analysis.explanation", classification.Explanation),
	)

	// Формируем результат анализа
	analysisResult := fmt.Sprintf(
		"Анализ наблюдения %s: Объект классифицирован как %s с уверенностью %.2f. %s",
		uuid,
		classification.ObjectType,
		classification.Confidence,
		classification.Explanation,
	)

	return analysisResult, classification.ObjectType, classification.Confidence, nil
}
