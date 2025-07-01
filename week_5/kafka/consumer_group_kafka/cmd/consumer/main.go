package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

const (
	brokerAddress = "localhost:9092"
	groupID       = "consumer-group"
	topicName     = "test-topic"
)

func main() {
	// Парсим параметры командной строки
	consumerID := flag.String("id", "consumer-1", "Consumer ID for identification in logs")
	flag.Parse()

	keepRunning := true
	log.Printf("[%s] starting Sarama consumer", *consumerID)

	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready: make(chan bool),
		id:    *consumerID,
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokerAddress, ","), groupID, config)
	if err != nil {
		log.Fatalf("[%s] failed to create consumer group client: %v", *consumerID, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		consume(ctx, client, consumer)
	}()

	<-consumer.ready
	log.Printf("[%s] consumer up and running, listening for messages on topic: %s", *consumerID, topicName)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Printf("[%s] terminating: context cancelled", *consumerID)
			keepRunning = false
		case <-sigterm:
			log.Printf("[%s] terminating: via signal", *consumerID)
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Fatalf("[%s] failed to close consumer group client: %v", *consumerID, err)
	}
}

func consume(ctx context.Context, client sarama.ConsumerGroup, consumer Consumer) {
	for {
		err := client.Consume(ctx, strings.Split(topicName, ","), &consumer)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			log.Fatalf("[%s] failed to consume: %v", consumer.id, err)
		}

		if ctx.Err() != nil {
			return
		}

		log.Printf("[%s] rebalancing", consumer.id)
		consumer.ready = make(chan bool)
	}
}
