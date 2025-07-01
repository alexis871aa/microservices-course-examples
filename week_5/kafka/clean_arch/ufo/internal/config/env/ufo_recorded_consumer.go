package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type ufoRecordedConsumerEnvConfig struct {
	Topic   string `env:"UFO_RECORDED_TOPIC_NAME,required"`
	GroupID string `env:"UFO_RECORDED_CONSUMER_GROUP_ID,required"`
}

type ufoRecordedConsumerConfig struct {
	raw ufoRecordedConsumerEnvConfig
}

func NewUfoRecordedConsumerConfig() (*ufoRecordedConsumerConfig, error) {
	var raw ufoRecordedConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &ufoRecordedConsumerConfig{raw: raw}, nil
}

func (cfg *ufoRecordedConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *ufoRecordedConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *ufoRecordedConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
