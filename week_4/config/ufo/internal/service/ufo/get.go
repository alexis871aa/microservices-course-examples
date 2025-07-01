package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	sighting, err := s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}

	return sighting, nil
}
