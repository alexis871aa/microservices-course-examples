package interceptor

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ufoMetrics "github.com/olezhek28/microservices-course-examples/week_7/metrics/ufo/internal/metrics"
)

// MetricsInterceptor создает gRPC интерцептор для записи метрик
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Засекаем время начала запроса
		start := time.Now()

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Записываем время выполнения
		duration := time.Since(start)
		durationSeconds := duration.Seconds()

		log.Printf("🕐 Request duration: %v (%f seconds) for method: %s", duration, durationSeconds, info.FullMethod)

		ufoMetrics.RequestDuration.Record(ctx, durationSeconds,
			metric.WithAttributes(
				attribute.String("method", info.FullMethod),
			),
		)

		// Определяем статус ответа
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		// Записываем метрику запроса
		statusLabel := "success"
		if statusCode != codes.OK {
			statusLabel = "error"
		}
		ufoMetrics.RequestsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", info.FullMethod),
				attribute.String("status", statusLabel),
			),
		)

		return resp, err
	}
}
