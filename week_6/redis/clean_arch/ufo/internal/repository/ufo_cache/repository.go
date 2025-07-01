package ufo_cache

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/cache"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/converter"
	repoModel "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/model"
)

const (
	cacheKeyPrefix = "ufo:sighting:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

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

	return repoConverter.SightingFromRedisView(ctx, sightingRedisView), nil
}

func (r *repository) Set(ctx context.Context, uuid string, sighting model.Sighting, ttl time.Duration) error {
	cacheKey := r.getCacheKey(uuid)

	redisView := repoConverter.SightingToRedisView(ctx, sighting)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}

func (r *repository) Delete(ctx context.Context, uuid string) error {
	cacheKey := r.getCacheKey(uuid)
	return r.cache.Del(ctx, cacheKey)
}
