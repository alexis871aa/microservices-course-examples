package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type UFOGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type UfoRecordedProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type UfoRecordedConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
