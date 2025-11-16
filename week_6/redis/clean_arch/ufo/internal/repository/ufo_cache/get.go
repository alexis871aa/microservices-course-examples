package ufo_cache

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/converter"
	repoModel "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	cacheKey := r.getCacheKey(uuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Sighting{}, model.ErrSightingNotFound
		}
		return model.Sighting{}, err
	}

	if len(values) == 0 {
		return model.Sighting{}, model.ErrSightingNotFound
	}

	var sightingRedisView repoModel.SightingRedisView
	err = redigo.ScanStruct(values, &sightingRedisView)
	if err != nil {
		return model.Sighting{}, err
	}

	return repoConverter.SightingFromRedisView(sightingRedisView), nil
}
