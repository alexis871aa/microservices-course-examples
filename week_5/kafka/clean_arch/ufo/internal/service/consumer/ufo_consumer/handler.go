package ufo_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.ufoRecordedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode UFORecorded", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.UUID),
		zap.String("location", event.Location),
		zap.String("description", event.Description),
	)

	return nil
}
