package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/cache"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/cache/redis"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/closer"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/logger"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/api/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/config"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/ufo"
	ufoCacheRepository "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/ufo_cache"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/service"
	ufoService "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/service/ufo"
)

type diContainer struct {
	ufoV1API ufoV1.UFOServiceServer

	ufoService service.UFOService

	ufoRepository      repository.UFORepository
	ufoCacheRepository repository.UFOCacheRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
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
		d.ufoService = ufoService.NewService(
			d.PartRepository(ctx),
			d.CacheRepository(),
			config.AppConfig().Redis.CacheTTL(),
		)
	}

	return d.ufoService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.UFORepository {
	if d.ufoRepository == nil {
		d.ufoRepository = ufoRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.ufoRepository
}

func (d *diContainer) CacheRepository() repository.UFOCacheRepository {
	if d.ufoCacheRepository == nil {
		d.ufoCacheRepository = ufoCacheRepository.NewRepository(d.RedisClient())
	}

	return d.ufoCacheRepository
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

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
