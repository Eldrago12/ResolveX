package main

import (
	"context"
	"time"

	"github.com/Eldrago12/ResolveX/internal/server"
	"github.com/joho/godotenv"

	"github.com/Eldrago12/ResolveX/internal/resolver"

	"github.com/Eldrago12/ResolveX/internal/prefetch"

	"github.com/Eldrago12/ResolveX/internal/limiter"

	"github.com/Eldrago12/ResolveX/internal/cache"

	"github.com/Eldrago12/ResolveX/config"
)

func main() {
	// Load .env file
	_ = godotenv.Load() // no fatal here; ignore error

	cfg := config.LoadConfig()

	redisCache := cache.NewCache(cfg.RedisURL)
	defer redisCache.Close()

	dnsResolver := resolver.NewResolver(redisCache)
	rateLimiter := limiter.NewRateLimiter(redisCache, cfg.RateLimit, cfg.RateLimitTTL)
	prefetcher := prefetch.NewPrefetcher(dnsResolver, redisCache, cfg.PrefetchLimit, time.Duration(cfg.PrefetchInterval)*time.Second)

	srv := &server.Server{
		Resolver:   dnsResolver,
		Limiter:    rateLimiter,
		Prefetcher: prefetcher,
		Port:       cfg.ServerPort,
	}

	ctx := context.Background()
	srv.StartPrefetcher(ctx)
	srv.Run()
}
