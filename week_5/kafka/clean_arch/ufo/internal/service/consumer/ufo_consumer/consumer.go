package ufo_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/logger"
	kafkaConverter "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/converter/kafka"
	def "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	ufoRecordedConsumer kafka.Consumer
	ufoRecordedDecoder  kafkaConverter.UFORecordedDecoder
}

func NewService(ufoRecordedConsumer kafka.Consumer, ufoRecordedDecoder kafkaConverter.UFORecordedDecoder) *service {
	return &service{
		ufoRecordedConsumer: ufoRecordedConsumer,
		ufoRecordedDecoder:  ufoRecordedDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ufoRecordedConsumer service")

	err := s.ufoRecordedConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
