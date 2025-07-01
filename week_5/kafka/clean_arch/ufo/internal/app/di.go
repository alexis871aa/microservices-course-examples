package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/closer"
	wrappedKafka "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka/producer"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/logger"
	kafkaMiddleware "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/middleware/kafka"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/api/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/config"
	kafkaConverter "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/converter/kafka"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/converter/kafka/decoder"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/repository"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/repository/ufo"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service"
	ufoConsumer "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service/consumer/ufo_consumer"
	ufoProducer "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service/producer/ufo_producer"
	ufoService "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service/ufo"
)

type diContainer struct {
	ufoV1API ufoV1.UFOServiceServer

	ufoService         service.UFOService
	ufoProducerService service.UFOProducerService
	ufoConsumerService service.ConsumerService

	ufoRepository repository.UFORepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	consumerGroup       sarama.ConsumerGroup
	ufoRecordedConsumer wrappedKafka.Consumer

	ufoRecordedDecoder  kafkaConverter.UFORecordedDecoder
	syncProducer        sarama.SyncProducer
	ufoRecordedProducer wrappedKafka.Producer
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
		d.ufoService = ufoService.NewService(d.PartRepository(ctx), d.UfoProducerService())
	}

	return d.ufoService
}

func (d *diContainer) UfoProducerService() service.UFOProducerService {
	if d.ufoProducerService == nil {
		d.ufoProducerService = ufoProducer.NewService(d.UFORecordedProducer())
	}

	return d.ufoProducerService
}

func (d *diContainer) UfoConsumerService() service.ConsumerService {
	if d.ufoConsumerService == nil {
		d.ufoConsumerService = ufoConsumer.NewService(d.UFORecordedConsumer(), d.UFORecordedDecoder())
	}

	return d.ufoConsumerService
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

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().UfoRecordedConsumer.GroupID(),
			config.AppConfig().UfoRecordedConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) UFORecordedConsumer() wrappedKafka.Consumer {
	if d.ufoRecordedConsumer == nil {
		d.ufoRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().UfoRecordedConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.ufoRecordedConsumer
}

func (d *diContainer) UFORecordedDecoder() kafkaConverter.UFORecordedDecoder {
	if d.ufoRecordedDecoder == nil {
		d.ufoRecordedDecoder = decoder.NewUFORecordedDecoder()
	}

	return d.ufoRecordedDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().UfoRecordedProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) UFORecordedProducer() wrappedKafka.Producer {
	if d.ufoRecordedProducer == nil {
		d.ufoRecordedProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().UfoRecordedProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.ufoRecordedProducer
}
