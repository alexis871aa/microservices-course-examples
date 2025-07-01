package ufo

import (
	"context"
)

func (s *service) Delete(ctx context.Context, uuid string) error {
	err := s.ufoRepository.Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
