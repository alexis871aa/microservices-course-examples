package main

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	ready chan bool
	id    string
}

// Setup запускается в начале новой сессии до вызова ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	log.Printf("[%s] consumer ready", c.id)
	return nil
}

// Cleanup запускается в конце жизни сессии
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("[%s] consumer cleanup", c.id)
	return nil
}

// ConsumeClaim обрабатывает сообщения из Kafka
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Код ниже не стоит перемещать в горутину, так как ConsumeClaim
	// уже запускается в горутине, см.:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("[%s] message channel was closed", c.id)
				return nil
			}

			log.Printf("[%s] message received: %s (partition: %d, offset: %d, timestamp: %v)",
				c.id, string(message.Value), message.Partition, message.Offset, message.Timestamp)
			session.MarkMessage(message, "")

		// Должен вернуться, когда `session.Context()` завершен.
		// В противном случае возникнет `ErrRebalanceInProgress` или `read tcp <ip>:<port>: i/o timeout` при перебалансировке кафки. см.:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
