package analysis

import (
	"context"

	analysisV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/analysis/v1"
	grpcClient "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/client/grpc"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
)

type client struct {
	client analysisV1.AnalysisServiceClient
}

// NewClient создает новый клиент для Analysis сервиса
func NewClient(grpcClient analysisV1.AnalysisServiceClient) grpcClient.AnalysisClient {
	return &client{
		client: grpcClient,
	}
}

// AnalyzeSighting анализирует наблюдение НЛО
func (c *client) AnalyzeSighting(ctx context.Context, uuid string) (model.AnalysisResult, error) {
	resp, err := c.client.AnalyzeSighting(ctx, &analysisV1.AnalyzeSightingRequest{
		Uuid: uuid,
	})
	if err != nil {
		return model.AnalysisResult{}, err
	}

	return model.AnalysisResult{
		AnalysisResult:  resp.AnalysisResult,
		Classification:  resp.Classification,
		ConfidenceScore: resp.ConfidenceScore,
	}, nil
}
