package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/closer"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
	analysisV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/analysis/v1"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_7/tracing/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/api/ufo/v1"
	grpcClient "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/client/grpc"
	analysisClient "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/client/grpc/analysis"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/config"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/repository"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/repository/ufo"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/service"
	ufoService "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/service/ufo"
)

type diContainer struct {
	ufoV1API ufoV1.UFOServiceServer

	ufoService service.UFOService

	ufoRepository repository.UFORepository

	analysisClient grpcClient.AnalysisClient

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UfoV1API(ctx context.Context) ufoV1.UFOServiceServer {
	if d.ufoV1API == nil {
		d.ufoV1API = ufoV1API.NewAPI(d.PartService(ctx))
	}

	return d.ufoV1API
}

func (d *diContainer) PartService(ctx context.Context) service.UFOService {
	if d.ufoService == nil {
		d.ufoService = ufoService.NewService(d.PartRepository(ctx), d.AnalysisClient(ctx))
	}

	return d.ufoService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.UFORepository {
	if d.ufoRepository == nil {
		d.ufoRepository = ufoRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.ufoRepository
}

func (d *diContainer) AnalysisClient(ctx context.Context) grpcClient.AnalysisClient {
	if d.analysisClient == nil {
		// Создаем gRPC соединение с интерцепторами
		conn, err := grpc.NewClient(
			"localhost:50052",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("analysis-service")),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create analysis connection: %s\n", err.Error()))
		}

		// Добавляем закрытие соединения в closer
		closer.AddNamed("Analysis gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		// Создаем proto клиента
		protoClient := analysisV1.NewAnalysisServiceClient(conn)

		// Создаем обертку
		d.analysisClient = analysisClient.NewClient(protoClient)
	}

	return d.analysisClient
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
