package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/olezhek28/microservices-course-examples/week_7/metrics/platform/pkg/logger"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_7/metrics/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_7/metrics/ufo/internal/interceptor"
	ufoMetrics "github.com/olezhek28/microservices-course-examples/week_7/metrics/ufo/internal/metrics"
)

// ufoServer простая реализация UFO gRPC сервиса
type ufoServer struct {
	ufoV1.UnimplementedUFOServiceServer
}

// Create создает новое наблюдение НЛО
func (s *ufoServer) Create(ctx context.Context, req *ufoV1.CreateRequest) (*ufoV1.CreateResponse, error) {
	logger.Info(ctx, "Create UFO sighting called",
		zap.String("location", req.Info.Location),
		zap.String("description", req.Info.Description))

	// Записываем метрики
	ufoMetrics.RequestsTotal.WithLabelValues("Create", "success").Inc()
	ufoMetrics.SightingsTotal.Inc()

	// Возвращаем фиктивный UUID
	return &ufoV1.CreateResponse{
		Uuid: "12345678-1234-1234-1234-123456789abc",
	}, nil
}

// Get возвращает наблюдение НЛО по идентификатору
func (s *ufoServer) Get(ctx context.Context, req *ufoV1.GetRequest) (*ufoV1.GetResponse, error) {
	logger.Info(ctx, "Get UFO sighting called", zap.String("uuid", req.Uuid))

	// Записываем метрики
	ufoMetrics.RequestsTotal.WithLabelValues("Get", "success").Inc()

	// Возвращаем фиктивные данные
	return &ufoV1.GetResponse{
		Sighting: &ufoV1.Sighting{
			Uuid: req.Uuid,
			Info: &ufoV1.SightingInfo{
				ObservedAt:      timestamppb.Now(),
				Location:        "Москва",
				Description:     "Треугольный объект в небе",
				Color:           wrapperspb.String("красный"),
				Sound:           wrapperspb.Bool(false),
				DurationSeconds: wrapperspb.Int32(300),
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}, nil
}

// Update обновляет существующее наблюдение НЛО
func (s *ufoServer) Update(ctx context.Context, req *ufoV1.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info(ctx, "Update UFO sighting called", zap.String("uuid", req.Uuid))

	// Записываем метрики
	ufoMetrics.RequestsTotal.WithLabelValues("Update", "success").Inc()

	return &emptypb.Empty{}, nil
}

// Delete выполняет мягкое удаление наблюдения НЛО
func (s *ufoServer) Delete(ctx context.Context, req *ufoV1.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info(ctx, "Delete UFO sighting called", zap.String("uuid", req.Uuid))

	// Записываем метрики
	ufoMetrics.RequestsTotal.WithLabelValues("Delete", "success").Inc()

	return &emptypb.Empty{}, nil
}

// AnalyzeSighting анализирует наблюдение НЛО
func (s *ufoServer) AnalyzeSighting(ctx context.Context, req *ufoV1.AnalyzeSightingRequest) (*ufoV1.AnalyzeSightingResponse, error) {
	logger.Info(ctx, "AnalyzeSighting called", zap.String("uuid", req.Uuid))

	// Записываем метрики
	ufoMetrics.RequestsTotal.WithLabelValues("AnalyzeSighting", "success").Inc()
	ufoMetrics.AnalysisRequestsTotal.Inc()

	// Симулируем время анализа
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		ufoMetrics.AnalysisDuration.Observe(duration.Seconds())
	}()

	// Возвращаем фиктивный результат анализа
	return &ufoV1.AnalyzeSightingResponse{
		AnalysisResult:  "Объект классифицирован как неизвестный с уверенностью 0.75",
		Classification:  "unknown",
		ConfidenceScore: 0.75,
	}, nil
}

func main() {
	ctx := context.Background()

	// Инициализация логгера
	err := logger.Init("info", false)
	if err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}

	// Регистрация UFO метрик
	ufoMetrics.RegisterMetrics()

	// Создание gRPC сервера с интерцептором метрик
	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.MetricsInterceptor()),
	)

	// Регистрация UFO сервиса
	ufoV1.RegisterUFOServiceServer(server, &ufoServer{})

	// Включение reflection для grpcurl
	reflection.Register(server)

	// Создание listener
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	logger.Info(ctx, "🚀 UFO gRPC server starting on :50051")

	// Запуск сервера
	if err := server.Serve(listener); err != nil {
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}
