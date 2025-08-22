package ufo

import (
	"context"
	"fmt"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	_, err := s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}

	return model.Sighting{}, fmt.Errorf("failed to get sighting: %w", err)
}
