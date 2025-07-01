package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/config/env"
)

var appConfig *config

type config struct {
	Logger              LoggerConfig
	UFOGRPC             UFOGRPCConfig
	Mongo               MongoConfig
	Kafka               KafkaConfig
	UfoRecordedProducer UfoRecordedProducerConfig
	UfoRecordedConsumer UfoRecordedConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	ufoGRPCCfg, err := env.NewUFOGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	ufoRecordedProducerCfg, err := env.NewUfoRecordedProducerConfig()
	if err != nil {
		return err
	}

	ufoRecordedConsumerCfg, err := env.NewUfoRecordedConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:              loggerCfg,
		UFOGRPC:             ufoGRPCCfg,
		Mongo:               mongoCfg,
		Kafka:               kafkaCfg,
		UfoRecordedProducer: ufoRecordedProducerCfg,
		UfoRecordedConsumer: ufoRecordedConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
