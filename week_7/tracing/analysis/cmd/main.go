package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/analysis/internal/handler"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/analysis/internal/service"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
	analysisV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/analysis/v1"
)

const (
	grpcPort = ":50052"
)

type config struct{}

func (c *config) CollectorEndpoint() string { return "localhost:4317" }
func (c *config) ServiceName() string       { return "analysis-service" }
func (c *config) Environment() string       { return "development" }
func (c *config) ServiceVersion() string    { return "1.0.0" }

func main() {
	ctx := context.Background()

	// Инициализация трейсинга
	cfg := &config{}
	err := tracing.InitTracer(ctx, cfg)
	if err != nil {
		log.Printf("Failed to initialize tracing: %v\n", err)
		return
	}
	defer func() {
		if cerr := tracing.ShutdownTracer(ctx); cerr != nil {
			log.Printf("Failed to shutdown tracer: %v\n", cerr)
		}
	}()

	// Создание сервиса и хендлера
	analysisService, err := service.NewAnalysisService()
	if err != nil {
		log.Printf("Failed to create analysis service: %v\n", err)
		return
	}

	analysisHandler := handler.NewAnalysisHandler(analysisService)

	// Создание gRPC сервера с трейсинг интерцептором
	server := grpc.NewServer(
		grpc.UnaryInterceptor(tracing.UnaryServerInterceptor("analysis-service")),
	)

	// Регистрация сервиса
	analysisV1.RegisterAnalysisServiceServer(server, analysisHandler)

	// Запуск сервера
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		return
	}

	log.Printf("Analysis service listening on %s", grpcPort)
	err = server.Serve(lis)
	if err != nil {
		log.Printf("Failed to serve: %v\n", err)
	}
}
