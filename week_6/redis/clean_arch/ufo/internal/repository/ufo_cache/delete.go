package ufo_cache

import "context"

func (r *repository) Delete(ctx context.Context, uuid string) error {
	cacheKey := r.getCacheKey(uuid)
	return r.cache.Del(ctx, cacheKey)
}
