package resolver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Eldrago12/ResolveX/internal/cache"

	"golang.org/x/sync/singleflight"
)

type Resolver struct {
	cache      *cache.Cache
	queryGroup singleflight.Group
}

func NewResolver(cache *cache.Cache) *Resolver {
	return &Resolver{cache: cache}
}

func (r *Resolver) ResolveDomain(ctx context.Context, domain string) (string, error) {
	if ip, err := r.cache.Get(ctx, domain); err == nil {
		return ip, nil
	}

	ip, err, _ := r.queryGroup.Do(domain, func() (interface{}, error) {
		ips, err := net.LookupIP(domain)
		if err != nil || len(ips) == 0 {
			return "", fmt.Errorf("no IP found for domain: %v", err)
		}
		var result string
		for _, ip := range ips {
			if ip.To4() != nil {
				result = ip.String()
				break
			} else {
				result = ip.String()
			}
		}

		r.cache.Set(ctx, domain, result, 10*time.Minute)
		r.cache.TrackFrequency(ctx, domain)
		return result, nil
	})

	return ip.(string), err
}
