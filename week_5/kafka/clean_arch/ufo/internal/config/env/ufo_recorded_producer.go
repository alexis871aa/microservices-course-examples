package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type ufoRecordedProducerEnvConfig struct {
	TopicName string `env:"UFO_RECORDED_TOPIC_NAME,required"`
}

type ufoRecordedProducerConfig struct {
	raw ufoRecordedProducerEnvConfig
}

func NewUfoRecordedProducerConfig() (*ufoRecordedProducerConfig, error) {
	var raw ufoRecordedProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &ufoRecordedProducerConfig{raw: raw}, nil
}

func (cfg *ufoRecordedProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *ufoRecordedProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
