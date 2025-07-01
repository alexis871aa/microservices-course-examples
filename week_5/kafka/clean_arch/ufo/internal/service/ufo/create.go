package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	err = s.ufoProducerService.ProduceUFORecorded(ctx, model.UFORecordedEvent{
		UUID:        uuid,
		ObservedAt:  info.ObservedAt,
		Location:    info.Location,
		Description: info.Description,
	})
	if err != nil {
		return "", err
	}

	return uuid, nil
}
