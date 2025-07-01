package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ufoMetrics "github.com/olezhek28/microservices-course-examples/week_7/metrics/ufo/internal/metrics"
)

// MetricsInterceptor создает gRPC интерцептор для записи метрик
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Увеличиваем счетчик активных соединений
		ufoMetrics.ActiveConnections.Inc()
		defer ufoMetrics.ActiveConnections.Dec()

		// Засекаем время начала запроса
		start := time.Now()

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Записываем время выполнения
		duration := time.Since(start)
		ufoMetrics.RequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

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
		ufoMetrics.RequestsTotal.WithLabelValues(info.FullMethod, statusLabel).Inc()

		return resp, err
	}
}
