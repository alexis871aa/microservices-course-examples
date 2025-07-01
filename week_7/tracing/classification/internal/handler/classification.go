package handler

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/classification/internal/service"
	classificationV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/classification/v1"
)

// ClassificationHandler - gRPC хендлер для классификации
type ClassificationHandler struct {
	classificationV1.UnimplementedClassificationServiceServer
	service *service.ClassificationService
}

// NewClassificationHandler создает новый хендлер
func NewClassificationHandler(service *service.ClassificationService) *ClassificationHandler {
	return &ClassificationHandler{
		service: service,
	}
}

// ClassifyObject обрабатывает запрос на классификацию объекта
func (h *ClassificationHandler) ClassifyObject(ctx context.Context, req *classificationV1.ClassifyObjectRequest) (*classificationV1.ClassifyObjectResponse, error) {
	objectType, confidence, explanation, err := h.service.ClassifyObject(
		ctx,
		req.Description,
		req.Color,
		req.DurationSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &classificationV1.ClassifyObjectResponse{
		ObjectType:  objectType,
		Confidence:  confidence,
		Explanation: explanation,
	}, nil
}
