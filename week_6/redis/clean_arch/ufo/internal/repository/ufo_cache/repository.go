package ufo_cache

import (
	"fmt"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/platform/pkg/cache"
)

const (
	cacheKeyPrefix = "ufo:sighting:"
)

type repository struct {
	// нам надо воспринимать как коннект к редису
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

// api -> service -> cache repo -> redis client (обертка наша) -> redis
func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}
