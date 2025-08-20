package main

import (
	"context"
	"errors"
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
	groupID       = "easy-kafka-group"
	topicName     = "test-topic"
	consumerID    = "consumer-1"
)

type Consumer struct {
	ready chan bool
	id    string
}

func main() {
	keepRunning := true
	log.Printf("[%s] starting Sarama consumer", consumerID)

	config := sarama.NewConfig()
	// конфигурируем консьюмер
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready: make(chan bool),
		id:    consumerID,
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokerAddress, ","), groupID, config)
	if err != nil {
		log.Fatalf("[%s] failed to create consumer group client: %v", consumerID, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		consume(ctx, client, consumer)
	}()

	<-consumer.ready
	log.Printf("[%s] consumer up and running, listening for messages on topic: %s", consumerID, topicName)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Printf("[%s] terminating: context cancelled", consumerID)
			keepRunning = false
		case <-sigterm:
			log.Printf("[%s] terminating: via signal", consumerID)
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Fatalf("[%s] failed to close consumer group client: %v", consumerID, err)
	}
}

func consume(ctx context.Context, client sarama.ConsumerGroup, consumer Consumer) {
	// здесь немного сложнее, чем в продюсере, потому что нужно учитывать ребалансировку при добавлении к примеру ещё одного констьюмера в группу
	for {
		// у нашего консьюмера есть метод Consume, который может слушать сразу несколько топиков (нужно на это обратить внимание) и третьим параметров передаётся хендлер и это бинго - интерфейс!!! в котором есть свои методы
		err := client.Consume(ctx, strings.Split(topicName, ","), &consumer)
		// очень важно для понимания, что client.Consume зависает в тот момент, когда мы консьюмера запустили. При этом она может разлочиться в нескольких ситуациях
		// - когда какая-то фатал ошибка
		// - а может развиснуть без ошибок, в случае ребалансинга и тогда Consume прекратит работать и сообщения перестанут обрабатываться и
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
		// инициализируем этот канал
		consumer.ready = make(chan bool)
	}
}

// Setup запускается в начале новой сессии до вызова ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// закрываем канал ready (подготовка для инициализации консьюмера)
	close(c.ready)
	log.Printf("[%s] consumer ready", c.id)
	return nil
}

// Cleanup запускается в конце жизни сессии
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	// тут можно добавлять логику очистки, если будет нужна
	log.Printf("[%s] consumer cleanup", c.id)
	return nil
}

// ConsumeClaim обрабатывает сообщения из Kafka
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// в этом месте мы реально читаем сообщения из кафки
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("[%s] message channel was closed", c.id)
				return nil
			}

			log.Printf("[%s] message received: %s (partition: %d, offset: %d)",
				c.id, string(message.Value), message.Partition, message.Offset)
			// помечаем сообщение как прочитанное
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
