package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/olezhek28/microservices-course-examples/week_4/di/platform/pkg/closer"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/di/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/api/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/config"
	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/repository"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/repository/ufo"
	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/service"
	ufoService "github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/service/ufo"
)

type diContainer struct {
	ufoV1API ufoV1.UFOServiceServer

	ufoService service.UFOService

	ufoRepository repository.UFORepository

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
		d.ufoService = ufoService.NewService(d.PartRepository(ctx))
	}

	return d.ufoService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.UFORepository {
	if d.ufoRepository == nil {
		d.ufoRepository = ufoRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.ufoRepository
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
