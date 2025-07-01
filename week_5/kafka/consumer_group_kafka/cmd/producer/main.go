package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
)

const (
	brokerAddress = "localhost:9092"
	topicName     = "test-topic"
)

func main() {
	producer, err := newSyncProducer([]string{brokerAddress})
	if err != nil {
		log.Fatalf("failed to start producer: %v", err)
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %v", err)
		}
	}()

	message := gofakeit.StreetName()
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}

	log.Printf("sending message to topic %s: %s", topicName, message)

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message to Kafka: %v", err)
		return
	}

	log.Printf("message successfully sent to partition %d with offset %d", partition, offset)
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
