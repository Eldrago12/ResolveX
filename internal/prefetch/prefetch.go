package prefetch

import (
	"context"
	"log"
	"time"

	"github.com/Eldrago12/ResolveX/internal/resolver"

	"github.com/Eldrago12/ResolveX/internal/cache"
)

type Prefetcher struct {
	resolver *resolver.Resolver
	cache    *cache.Cache
	limit    int
	interval time.Duration
}

func NewPrefetcher(resolver *resolver.Resolver, cache *cache.Cache, limit int, interval time.Duration) *Prefetcher {
	return &Prefetcher{
		resolver: resolver,
		cache:    cache,
		limit:    limit,
		interval: interval,
	}
}

func (p *Prefetcher) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.prefetchTopDomains(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Prefetcher) prefetchTopDomains(ctx context.Context) {
	domains, err := p.cache.GetTopDomains(ctx, p.limit)
	if err != nil {
		log.Printf("Failed to retrieve top domains for prefetch: %v", err)
		return
	}

	for _, domain := range domains {
		go func(d string) {
			if _, err := p.resolver.ResolveDomain(ctx, d); err != nil {
				log.Printf("Prefetch failed for domain %s: %v", d, err)
			}
		}(domain)
	}
}
