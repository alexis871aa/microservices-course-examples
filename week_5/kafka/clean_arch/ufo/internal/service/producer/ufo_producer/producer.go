package ufo_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/kafka"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/platform/pkg/logger"
	eventsV1 "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/shared/pkg/proto/events/v1"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/model"
)

type service struct {
	ufoRecordedProducer kafka.Producer
}

func NewService(ufoRecordedProducer kafka.Producer) *service {
	return &service{
		ufoRecordedProducer: ufoRecordedProducer,
	}
}

func (p *service) ProduceUFORecorded(ctx context.Context, event model.UFORecordedEvent) error {
	var observedAt *timestamppb.Timestamp
	if event.ObservedAt != nil {
		observedAt = timestamppb.New(*event.ObservedAt)
	}

	msg := &eventsV1.UFORecorded{
		ObservedAt:  observedAt,
		Location:    event.Location,
		Description: event.Description,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	err = p.ufoRecordedProducer.Send(ctx, []byte(event.UUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish UFORecorded", zap.Error(err))
		return err
	}

	return nil
}
