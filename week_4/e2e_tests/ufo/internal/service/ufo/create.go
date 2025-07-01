package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
