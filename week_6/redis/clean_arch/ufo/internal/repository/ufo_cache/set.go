package ufo_cache

import (
	"context"
	"time"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/converter"
)

func (r *repository) Set(ctx context.Context, uuid string, sighting model.Sighting, ttl time.Duration) error {
	cacheKey := r.getCacheKey(uuid)

	redisView := repoConverter.SightingToRedisView(sighting)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
