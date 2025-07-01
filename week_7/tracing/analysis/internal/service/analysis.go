package service

import (
	"context"
	"fmt"

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
	// Создаем спан для получения данных о наблюдении
	ctx, span := tracing.StartSpan(ctx, "analysis.get_sighting")

	// Получаем данные о наблюдении из UFO сервиса
	sighting, err := s.ufoClient.Get(ctx, &ufoV1.GetRequest{Uuid: uuid})
	if err != nil {
		span.End()
		return "", "", 0, fmt.Errorf("failed to get sighting: %w", err)
	}
	span.End()

	// Создаем спан для классификации
	ctx, span = tracing.StartSpan(ctx, "analysis.classify")
	defer span.End()

	// Отправляем данные в Classification сервис
	classification, err := s.classificationClient.ClassifyObject(ctx, &classificationV1.ClassifyObjectRequest{
		Description:     sighting.Sighting.Info.Description,
		Color:           sighting.Sighting.Info.Color.GetValue(),
		DurationSeconds: sighting.Sighting.Info.DurationSeconds.GetValue(),
	})
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to classify object: %w", err)
	}

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
