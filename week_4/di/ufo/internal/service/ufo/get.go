package ufo

import (
	"context"
	"github.com/olezhek28/microservices-course-examples/week_4/di/platform/pkg/logger"
	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/model"
	"go.uber.org/zap"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	sighting, err := s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "failed to get ufo",
			zap.String("uuid", uuid),
			zap.Error(err),
		)
		return model.Sighting{}, err
	}

	return sighting, nil
}
