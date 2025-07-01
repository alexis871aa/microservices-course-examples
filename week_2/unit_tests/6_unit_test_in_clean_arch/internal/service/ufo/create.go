package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
