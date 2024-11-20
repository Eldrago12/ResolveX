package limiter

import (
	"context"
	"log"
	"time"

	"github.com/Eldrago12/ResolveX/internal/cache"
)

type RateLimiter struct {
	cache        *cache.Cache
	limit        int
	rateLimitTTL time.Duration
}

func NewRateLimiter(cache *cache.Cache, limit int, ttl int) *RateLimiter {
	return &RateLimiter{
		cache:        cache,
		limit:        limit,
		rateLimitTTL: time.Duration(ttl) * time.Second,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	count, err := rl.cache.Increment(ctx, key, rl.rateLimitTTL)
	if err != nil {
		log.Printf("Error with rate limiter increment for key %s: %v", key, err)
		return false, err
	}

	log.Printf("Request count for %s: %d (limit: %d)", key, count, rl.limit)

	return count <= int64(rl.limit), nil
}
